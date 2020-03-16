use actix_identity::Identity;
use actix_web::{get, http, post, web, HttpResponse};
use askama::Template;
use serde::Deserialize;

use crate::templates;
use crate::DbPool;
use crate::{redirect_to, serve_template, unwrap_or_redirect};

pub mod passwords;
pub mod permissions;
pub mod users;

#[derive(Deserialize)]
struct LoginForm {
    username: String,
    password: String,
    target: String,
}

#[derive(Deserialize)]
struct LoginQuery {
    error: Option<String>,
}

#[get("/login")]
async fn login(info: web::Query<LoginQuery>) -> HttpResponse {
    serve_template!(templates::Login {
        header_links: &[],
        target: "/",
        error: info.error.is_some(),
    });
}

#[post("/login")]
async fn login_form(
    id: Identity,
    pool: web::Data<DbPool>,
    form: web::Form<LoginForm>,
) -> HttpResponse {
    // we store this here since `form` moves into the closure
    let target = form.target.clone();

    let conn = pool.get().expect("Error getting db connection from pool");
    let user = unwrap_or_redirect!({
        web::block(move || {
            users::check_login(form.username.to_owned(), form.password.to_owned(), &conn)
        }).await
    } => "login?error=1");

    id.remember(serde_json::to_string(&user).unwrap());

    redirect_to!(target);
}

#[derive(Deserialize)]
struct SignupForm {
    username: String,
    password: String,
    target: String,
}

#[get("/signup")]
async fn signup() -> HttpResponse {
    serve_template!(templates::Signup {
        header_links: &[],
        target: "/",
        error: false,
    });
}

#[post("/signup")]
async fn signup_form(
    id: Identity,
    pool: web::Data<DbPool>,
    form: web::Form<SignupForm>,
) -> HttpResponse {
    let target = form.target.clone();

    let conn = pool.get().expect("Error getting database");
    let user = unwrap_or_redirect!({
        web::block(move || {
            users::create_user(form.username.to_owned(), form.password.to_owned(), &conn)
        }).await
    } => "/signup?error=1");

    id.remember(serde_json::to_string(&user).unwrap());

    redirect_to!(target);
}
