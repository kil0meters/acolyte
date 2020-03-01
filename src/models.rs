use crate::auth::permissions::AuthLevel;
use diesel::{Insertable, Queryable};
use std::time;

use crate::schema::accounts;

#[derive(Insertable, Queryable, Identifiable, Debug, PartialEq)]
#[table_name = "accounts"]
pub struct Account {
    pub id: String,
    pub username: String,
    pub password_hash: String,
    pub created_at: time::SystemTime,
    pub permissions: AuthLevel,
}
