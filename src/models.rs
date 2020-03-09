use std::str::FromStr;
use std::time;

use anyhow::{anyhow, Result};
use chrono::prelude::*;
use diesel::{Insertable, Queryable};
use rand::distributions::Alphanumeric;
use rand::Rng;
use regex::Regex;
use serde::{Deserialize, Serialize};

use crate::auth::{passwords, permissions};
use crate::schema::{accounts, blog_posts, posts};

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

/// Struct containing account data
/// ````
/// struct Account {
///     id: String,
///     username: String,
///     password_hash: String, // not serialized
///     created_at: time::SystemTime,
///     permissions: permissions::AuthLevel,
/// }
/// ```
#[derive(Serialize, Deserialize, Insertable, Queryable, Identifiable, Debug, PartialEq)]
#[table_name = "accounts"]
pub struct Account {
    pub id: String,

    pub username: String,

    #[serde(skip_serializing, default = "string_default")]
    pub password_hash: String,

    pub created_at: time::SystemTime,

    pub permissions: permissions::AuthLevel,
}

impl Account {
    /// Creates a new account object, hashing the password
    pub fn new(username: String, password: String) -> Account {
        let id = create_id(ID_LENGTH);

        // TODO: handle error
        let password_hash = passwords::hash_password(password).unwrap();

        Account {
            id,
            username: username.to_lowercase().trim().to_owned(),
            password_hash,
            created_at: time::SystemTime::now(),
            permissions: permissions::STANDARD,
        }
    }

    /// Checks an account ot see if it's valid
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

/// Struct containing a post
/// ```
/// struct Post {
///     id: String,
///     account_id: String,
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
#[belongs_to(Account)]
#[table_name = "posts"]
pub struct Post {
    pub id: String,
    pub account_id: String,
    pub title: String,
    pub body: Option<String>,
    pub link: Option<String>,
    pub removed: bool,
    pub created_at: time::SystemTime,

    // this should possibly be stored as some sort of BigInt
    // but that's an issue for another time
    pub upvotes: i32,
    pub downvotes: i32,
}

impl Post {
    pub fn new(account_id: String, title: String, body: String, link: String) -> Post {
        let id = create_id(ID_LENGTH);

        let body: Option<String> = if body.trim().is_empty() {
            None
        } else {
            Some(body)
        };

        let link: Option<String> = if link.trim().is_empty() {
            None
        } else {
            Some(link)
        };

        Post {
            id,
            account_id,
            title,
            body,
            link,
            removed: false,
            created_at: time::SystemTime::now(),
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
    pub last_modified: chrono::NaiveDateTime,
    pub created_at: chrono::NaiveDateTime,
}

impl FromStr for BlogPost {
    type Err = anyhow::Error;

    fn from_str(post: &str) -> Result<BlogPost> {
        lazy_static! {
            /// Matching group 1: inside yaml metadata
            ///
            /// Matching group 2: Post content
            static ref METADATA_REGEX: Regex =
                Regex::new(r"---\n((?:.|\n)*)---\n((?:.|\n)*)").unwrap();

            /// Mathces everything that isn't alphanumeric
            static ref NON_ALPHANUM_REGEX: Regex =
                Regex::new(r"[^A-Za-z0-9 ]").unwrap();
        }

        if let Some(caps) = METADATA_REGEX.captures(post.trim()) {
            let metadata_str = caps.get(1).map_or("", |m| m.as_str());
            let post_content = caps.get(2).map_or("", |m| m.as_str());

            #[derive(Deserialize)]
            struct PostMetadata {
                title: String,
                date: String,
            }

            let metadata = serde_yaml::from_str::<PostMetadata>(&metadata_str)?;

            let date = NaiveDate::parse_from_str(&metadata.date, "%Y-%m-%d")?.and_hms(0, 0, 0);

            let id = NON_ALPHANUM_REGEX
                .replace(&metadata.title, "")
                .to_lowercase()
                .replace(" ", "_");

            Ok(BlogPost {
                id,
                title: metadata.title,
                body: post_content.to_owned(),
                last_modified: date,
                created_at: date,
            })
        } else {
            debug!(
                "Encountered error matching the following post data:\n\"{}\"",
                post.trim()
            );
            Err(anyhow!("Invalid post metadata"))
        }
    }
}
