use diesel::pg::PgConnection;
use diesel::prelude::*;
use diesel::result::Error;

use crate::auth::passwords;
use crate::models;
use crate::schema;

/// Checks whether a username/password pair is valid.
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

    if let Some(account) = account {
        let passwords_match =
            passwords::verify_password(password, account.password_hash.clone()).unwrap();

        if passwords_match {
            Ok(Some(account))
        } else {
            Ok(None)
        }
    } else {
        Ok(None)
    }
}

/// Creates a `models::Account` object and adds it to the database.
pub fn create_account(
    username: String,
    password: String,
    conn: &PgConnection,
) -> Result<Option<models::Account>, Error> {
    let account = models::Account::new(username, password);

    if account.validate() {
        let result = diesel::insert_into(schema::accounts::table)
            .values(&account)
            .execute(conn)?;
        println!("result: {}", result);

        Ok(Some(account))
    } else {
        // Account is invalid
        Ok(None)
    }
}
