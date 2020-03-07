use std::time;

use diesel::{Insertable, Queryable};
use rand::distributions::Alphanumeric;
use rand::Rng;
use regex::Regex;
use serde::{Deserialize, Serialize};

use crate::auth::{passwords, permissions};
use crate::schema::{accounts, posts};

fn string_default() -> String {
    "".to_owned()
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
        let id = rand::thread_rng()
            .sample_iter(&Alphanumeric)
            .take(7)
            .collect::<String>();

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

    /// Validates an account ot see if it's valid
    pub fn validate(&self) -> bool {
        // avoid compiling the regex every time
        lazy_static! {
            static ref USERNAME_REGEX: Regex = Regex::new(r"^[a-zA-Z][a-zA-Z0-9_]{0,15}$").unwrap();
        }

        USERNAME_REGEX.is_match(&self.username)
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
    // but that's an issue for another tmie
    pub upvotes: i32,
    pub downvotes: i32,
}
