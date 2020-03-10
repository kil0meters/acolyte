use actix_identity::Identity;
use actix_web::{error, get, http, post, web, Error, HttpRequest, HttpResponse};
use askama::Template;
use serde::{Deserialize, Serialize};
use serde_json::json;

use crate::models;
use crate::templates;
use crate::DbPool;

pub mod passwords;
pub mod permissions;
pub mod users;

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
    let user = web::block(move || {
        users::check_login(form.username.to_owned(), form.password.to_owned(), &conn)
    })
    .await
    .map_err(|_| HttpResponse::InternalServerError())?;

    if let Some(user) = user {
        id.remember(serde_json::to_string(&user).unwrap());

        Ok(HttpResponse::SeeOther()
            .header(http::header::LOCATION, target.to_owned())
            .finish())
    } else {
        // If the user didn't exist, redirect with the error optioni
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

    let user = web::block(move || {
        users::create_user(form.username.to_owned(), form.password.to_owned(), &conn)
    })
    .await;

    println!("result: {:?}", user);

    match user {
        Ok(user) => {
            id.remember(serde_json::to_string(&user).unwrap());

            HttpResponse::Found()
                .header(http::header::LOCATION, target.to_owned())
                .finish()
        }
        Err(_) => HttpResponse::Found()
            .header(http::header::LOCATION, "/signup?error=1")
            .finish(),
    }
}
