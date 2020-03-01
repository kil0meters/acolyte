use diesel::pg::PgConnection;
use diesel::prelude::*;
use diesel::result::Error;

use crate::models;

pub fn check_login(
    username: String,
    password: String,
    conn: &PgConnection,
) -> Result<Option<models::Account>, Error> {
    use crate::schema::accounts::dsl;

    let account = dsl::accounts
        .filter(dsl::username.eq(username))
        .first::<models::Account>(conn)
        .optional()?;

    Ok(account)
}

// fn create_account() -> Result<Option<models::Account>, Error> {}
