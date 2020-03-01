use actix_identity::Identity;
use actix_web::{error, get, http, post, web, Error, HttpRequest, HttpResponse};
use serde::{Deserialize, Serialize};
use serde_json::json;
use tera;

use crate::DbPool;

pub mod accounts;
pub mod permissions;

#[derive(Serialize, Deserialize)]
pub struct Account {
    pub username: String,
    pub password: String,
}

#[derive(Deserialize)]
struct LoginForm {
    username: String,
    password: String,
    target: String,
}

#[get("/login")]
async fn login(id: Identity, tmpl: web::Data<tera::Tera>) -> Result<HttpResponse, Error> {
    let ctx = tera::Context::from_value(json!({
        "title": "Login - milesbenton.com",
        "page_title": "Login",
        "header": [],
        "target": "/",
        "error": false,
    }))
    .unwrap();

    let s = tmpl
        .render("login.html", &ctx)
        .map_err(|e| error::ErrorInternalServerError(format!("Template error: {}", e)))?;

    Ok(HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(s))
}

#[post("/login")]
async fn login_form(
    pool: web::Data<DbPool>,
    form: web::Form<LoginForm>,
    id: Identity,
) -> Result<HttpResponse, Error> {
    let conn = pool.get().expect("Error getting db connection from pool");
    let account = accounts::check_login(form.username.to_owned(), form.password.to_owned(), &conn)
        .expect("Encountered an error when trying to read account");

    if let Some(account) = account {
        // web::block(|| {
        // })

        id.remember(
            serde_json::to_string(&Account {
                username: form.username.to_owned(),
                password: form.password.to_owned(),
            })
            .unwrap(),
        );

        Ok(HttpResponse::Found()
            .header(http::header::LOCATION, form.target.to_owned())
            .finish())
    } else {
        // If the account didn't exist, redirect with the error optioni
        Ok(HttpResponse::Found()
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
async fn signup(tmpl: web::Data<tera::Tera>) -> Result<HttpResponse, Error> {
    let ctx = tera::Context::from_value(json!({
        "title": "Signup - milesbenton.com",
        "page_title": "Signup",
        "header": [],
        "target": "/",
        "error": false,
    }))
    .unwrap();

    let s = tmpl
        .render("signup.html", &ctx)
        .map_err(|e| error::ErrorInternalServerError(format!("Template error: {}", e)))?;

    Ok(HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(s))
}

#[post("/signup")]
async fn signup_form(form: web::Form<SignupForm>, id: Identity) -> Result<HttpResponse, Error> {
    id.remember(
        serde_json::to_string(&Account {
            username: form.username.to_owned(),
            password: form.password.to_owned(),
        })
        .unwrap(),
    );

    id.remember(form.username.to_owned());
    id.remember(form.password.to_owned());

    Ok(HttpResponse::Found()
        .header(http::header::LOCATION, form.target.to_owned())
        .finish())
}
