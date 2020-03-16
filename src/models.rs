use std::str::FromStr;

use anyhow::{anyhow, Result};
use chrono::prelude::*;
use diesel::{Insertable, Queryable};
use rand::distributions::Alphanumeric;
use rand::Rng;
use regex::Regex;
use serde::{Deserialize, Serialize};

use crate::auth::{passwords, permissions};
use crate::schema::{blog_posts, comments, threads, users};

const ID_LENGTH: usize = 8;

fn string_default() -> String {
    "".to_owned()
}

fn create_id(length: usize) -> String {
    // TODO: with size of 8, this gives a 1 in 62^8 chance of an error happening
    rand::thread_rng()
        .sample_iter(&Alphanumeric)
        .take(length)
        .collect::<String>()
}

/// Struct containing user data
/// ````
/// struct User {
///     id: String,
///     username: String,
///     password_hash: String, // not serialized
///     updated_at: chrono::NaiveDateTime,
///     created_at: time::SystemTime,
///     permissions: permissions::AuthLevel,
/// }
/// ```
#[derive(Serialize, Deserialize, Insertable, Queryable, Identifiable, Clone, Debug, PartialEq)]
#[table_name = "users"]
pub struct User {
    pub id: String,
    pub username: String,
    #[serde(skip_serializing, default = "string_default")]
    pub password_hash: String,
    pub updated_at: chrono::NaiveDateTime,
    pub created_at: chrono::NaiveDateTime,
    pub permissions: permissions::AuthLevel,
}

impl User {
    /// Creates a new user object, hashing the password
    pub fn new(username: String, password: String) -> User {
        let id = create_id(ID_LENGTH);

        // TODO: handle error
        let password_hash = passwords::hash_password(password).unwrap();
        let now = Utc::now().naive_utc();

        User {
            id,
            username: username.to_lowercase().trim().to_owned(),
            password_hash,
            created_at: now,
            updated_at: now,
            permissions: permissions::STANDARD,
        }
    }

    pub fn from_identity(id: actix_identity::Identity) -> User {
        lazy_static! {
            static ref DEFAULT_USER: User = User {
                id: String::new(),
                username: "ANON".to_owned(),
                password_hash: String::new(),
                created_at: NaiveDateTime::from_timestamp(0, 0),
                updated_at: NaiveDateTime::from_timestamp(0, 0),
                permissions: permissions::LOGGED_OUT,
            };
        }

        if let Some(id) = id.identity() {
            if let Ok(user) = serde_json::from_str::<User>(&id) {
                debug!("Identity: {:?}", user);
                return user;
            }
        }

        DEFAULT_USER.clone()
    }

    /// Checks an user ot see if it's valid
    pub fn validate(&self) -> Result<()> {
        // avoid compiling the regex every time
        lazy_static! {
            static ref USERNAME_REGEX: Regex = Regex::new(r"^[a-zA-Z][a-zA-Z0-9_]{0,15}$").unwrap();
        }

        if USERNAME_REGEX.is_match(&self.username) {
            Ok(())
        } else {
            Err(anyhow!("Invalid username"))
        }
    }
}

/// Struct containing a forum thread
/// ```
/// struct Thread {
///     id: String,
///     user_id: String,
///     username: String,
///     title: String,
///     body: Option<String>,
///     link: Option<String>,
///     removed: bool,
///     created_at: time::SystemTime,
///     upvotes: i32,
///     downvotes: i32,
/// }
/// ```
#[derive(Associations, Insertable, Queryable, Debug, PartialEq)]
#[belongs_to(User)]
#[table_name = "threads"]
pub struct Thread {
    pub id: String,
    pub user_id: String,
    pub username: String,
    pub title: String,
    pub body: Option<String>,
    pub body_html: Option<String>,
    pub link: Option<String>,
    pub removed: bool,
    pub updated_at: chrono::NaiveDateTime,
    pub created_at: chrono::NaiveDateTime,

    // this should possibly be stored as some sort of BigInt
    // but that's an issue for anoher time
    pub upvotes: i32,
    pub downvotes: i32,
}

impl Thread {
    // TODO: Sanitize html
    pub fn new(author: User, title: String, body: String, link: String) -> Thread {
        use pulldown_cmark::{html, Options, Parser};

        let id = create_id(ID_LENGTH);

        let (body, body_html) = if body.trim().is_empty() {
            (None, None)
        } else {
            let parser = Parser::new_ext(&body, Options::all());

            let mut unsafe_body_html = String::new();
            html::push_html(&mut unsafe_body_html, parser);

            let body_html = ammonia::clean(&*unsafe_body_html);

            (Some(body), Some(body_html))
        };

        let link: Option<String> = if link.trim().is_empty() {
            None
        } else {
            Some(link)
        };

        let now = Utc::now().naive_utc();

        Thread {
            id,
            user_id: author.id,
            username: author.username,
            title,
            body,
            body_html,
            link,
            removed: false,
            updated_at: now,
            created_at: now,
            upvotes: 0,
            downvotes: 0,
        }
    }

    pub fn validate(&self) -> Result<()> {
        lazy_static! {
            static ref URL_REGEX: Regex =
                Regex::new(r"^http://[a-zA-Z0-9_\-]+\.[a-zA-Z0-9_\-]+\.[a-zA-Z0-9_\-]+$").unwrap();
        }

        // TODO: there has to be a better way to write this
        if let Some(link) = &self.link {
            if URL_REGEX.is_match(&link) {
                Err(anyhow!("URL {} is invalid", link))
            } else {
                Ok(())
            }
        } else {
            Ok(())
        }
    }
}

#[derive(Insertable, Queryable, Debug, PartialEq)]
#[table_name = "blog_posts"]
pub struct BlogPost {
    pub id: String,
    pub title: String,
    pub body: String,
    pub updated_at: chrono::NaiveDateTime,
    pub created_at: chrono::NaiveDateTime,
}

impl FromStr for BlogPost {
    type Err = anyhow::Error;

    fn from_str(post: &str) -> Result<BlogPost> {
        lazy_static! {
            // \r\n has to be used because of some weird meme with how forms are escaped
            // https://github.com/rust-lang/regex/issues/244
            static ref METADATA_REGEX: Regex =
                Regex::new(r"(?m)---\r\n((?:.|\n)*)---\r\n((?:.|\n)*)").unwrap();

            /// Mathces everything that isn't alphanumeric
            static ref NON_ALPHANUM_REGEX: Regex =
                Regex::new(r"[^A-Za-z0-9 ]").unwrap();
        }

        if let Some(caps) = METADATA_REGEX.captures(post.trim()) {
            let metadata_str = caps.get(1).map_or("", |m| m.as_str());
            let post_content = caps.get(2).map_or("", |m| m.as_str());

            #[derive(Deserialize)]
            struct ThreadMetadata {
                title: String,
                date: String,
            }

            debug!("Metadata:\n\"{}\"", metadata_str);
            let metadata = serde_yaml::from_str::<ThreadMetadata>(&metadata_str)?;

            let date = NaiveDate::parse_from_str(&metadata.date, "%Y-%m-%d")?.and_hms(0, 0, 0);

            let id = NON_ALPHANUM_REGEX
                .replace(&metadata.title, "")
                .to_lowercase()
                .replace(" ", "_");

            Ok(BlogPost {
                id,
                title: metadata.title,
                body: post_content.to_owned(),
                updated_at: date,
                created_at: date,
            })
        } else {
            Err(anyhow!("Invalid post data:\n\"{}\"", post.trim()))
        }
    }
}

#[derive(Insertable, Queryable, Debug, Clone, PartialEq)]
#[table_name = "comments"]
pub struct Comment {
    pub id: String,
    pub id_parents: String,
    pub user_id: String,
    pub username: String,
    pub body: String,
    pub body_html: String,
    pub removed: bool,
    pub updated_at: chrono::NaiveDateTime,
    pub created_at: chrono::NaiveDateTime,
    pub upvotes: i32,
    pub downvotes: i32,
    // pub chilidren: Vec<Comment>,
}

impl Comment {
    pub fn new(author: User, parent_id: String, body: String) -> Comment {
        use pulldown_cmark::{html, Options, Parser};

        let parser = Parser::new_ext(&body, Options::all());

        let mut unsafe_body_html = String::new();
        html::push_html(&mut unsafe_body_html, parser);

        let body_html = ammonia::clean(&*unsafe_body_html);

        let now = Utc::now().naive_utc();

        let id = create_id(ID_LENGTH);

        Comment {
            // poor ordering but this is more efficient
            id_parents: format!("{}-{}", parent_id, id),
            id,
            user_id: author.id,
            username: author.username,
            body,
            body_html,
            removed: false,
            updated_at: now,
            created_at: now,
            upvotes: 0,
            downvotes: 0,
        }
    }

    pub fn validate(&self) -> Result<()> {
        Ok(())
    }
}
