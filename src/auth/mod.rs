use actix_identity::Identity;
use actix_web::{error, get, http, post, web, Error, HttpRequest, HttpResponse};
use askama::Template;
use serde::{Deserialize, Serialize};
use serde_json::json;

use crate::models;
use crate::templates;
use crate::DbPool;

pub mod accounts;
pub mod passwords;
pub mod permissions;

#[derive(Deserialize)]
struct LoginForm {
    username: String,
    password: String,
    target: String,
}

#[get("/login")]
async fn login(id: Identity) -> Result<HttpResponse, Error> {
    let s = templates::Login {
        header_links: &[],
        target: "/",
        error: false,
    }
    .render()
    .unwrap();

    Ok(HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(s))
}

#[post("/login")]
async fn login_form(
    id: Identity,
    pool: web::Data<DbPool>,
    form: web::Form<LoginForm>,
) -> Result<HttpResponse, Error> {
    let conn = pool.get().expect("Error getting db connection from pool");

    // we store this here since `form` moves into the closure
    let target = form.target.clone();

    // Don't block server thread with DB query
    let account = web::block(move || {
        accounts::check_login(form.username.to_owned(), form.password.to_owned(), &conn)
    })
    .await
    .map_err(|_| HttpResponse::InternalServerError().finish())?;

    if let Some(account) = account {
        id.remember(serde_json::to_string(&account).unwrap());

        Ok(HttpResponse::SeeOther()
            .header(http::header::LOCATION, target.to_owned())
            .finish())
    } else {
        // If the account didn't exist, redirect with the error optioni
        Ok(HttpResponse::SeeOther()
            .header(http::header::LOCATION, "/login?error=1")
            .finish())
    }
}

#[derive(Deserialize)]
struct SignupForm {
    username: String,
    password: String,
    target: String,
}

#[get("/signup")]
async fn signup() -> Result<HttpResponse, Error> {
    let s = templates::Signup {
        header_links: &[],
        target: "/",
        error: false,
    }
    .render()
    .unwrap();

    Ok(HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(s))
}

#[post("/signup")]
async fn signup_form(
    id: Identity,
    pool: web::Data<DbPool>,
    form: web::Form<SignupForm>,
) -> HttpResponse {
    let conn = pool.get().expect("Error getting database");
    // same as above
    let target = form.target.clone();

    let account = web::block(move || {
        accounts::create_account(form.username.to_owned(), form.password.to_owned(), &conn)
    })
    .await;

    println!("result: {:?}", account);

    match account {
        Ok(account) => {
            id.remember(serde_json::to_string(&account).unwrap());

            HttpResponse::Found()
                .header(http::header::LOCATION, target.to_owned())
                .finish()
        }
        Err(_) => HttpResponse::Found()
            .header(http::header::LOCATION, "/signup?error=1")
            .finish(),
    }
}
