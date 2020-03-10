use anyhow::Result;
use diesel::pg::PgConnection;
use diesel::prelude::*;

use crate::auth::passwords;
use crate::models::User;
use crate::schema;

/// Checks whether a username/password pair is valid.
pub fn check_login(
    username: String,
    password: String,
    conn: &PgConnection,
) -> Result<Option<User>> {
    use crate::schema::users::dsl;

    let user = dsl::users
        .filter(dsl::username.eq(username))
        .first::<User>(conn)
        .optional()?;

    if let Some(user) = user {
        let passwords_match =
            passwords::verify_password(password, user.password_hash.clone()).unwrap();

        if passwords_match {
            Ok(Some(user))
        } else {
            Ok(None)
        }
    } else {
        Ok(None)
    }
}

/// Creates a `User` object and adds it to the database.
pub fn create_user(username: String, password: String, conn: &PgConnection) -> Result<User> {
    let user = User::new(username, password);

    user.validate()?;

    let result = diesel::insert_into(schema::users::table)
        .values(&user)
        .execute(conn)?;

    println!("result: {}", result);

    Ok(user)
}
