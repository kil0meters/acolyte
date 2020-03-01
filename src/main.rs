#![warn(clippy::all)]

#[macro_use]
extern crate diesel;

use std::{env, time};

use actix::Actor;
use actix_identity::{CookieIdentityPolicy, IdentityService};
use actix_web::{middleware, web, App, HttpServer};
use diesel::pg::PgConnection;
use diesel::r2d2::{ConnectionManager, Pool};
use tera::Tera;

pub mod auth;
pub mod chat;
pub mod frontpage;
pub mod models;
pub mod schema;

type DbPool = Pool<ConnectionManager<PgConnection>>;

#[actix_rt::main]
async fn main() -> std::io::Result<()> {
    println!("Serving on 127.0.0.1:8080");

    std::env::set_var("RUST_LOG", "actix_web=info");
    env_logger::init();

    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL is required.");
    let manager = ConnectionManager::<PgConnection>::new(database_url);
    let pool = Pool::builder()
        .build(manager)
        .expect("Failed to create database pool");

    let srv = chat::server::Server::default().start();

    HttpServer::new(move || {
        let mut tera = Tera::default();
        tera.add_raw_templates(vec![
            (
                "base.html",
                include_str!(concat!(env!("CARGO_MANIFEST_DIR"), "/templates/base.html")),
            ),
            (
                "static_header.html",
                include_str!(concat!(
                    env!("CARGO_MANIFEST_DIR"),
                    "/templates/static_header.html"
                )),
            ),
            (
                "home.html",
                include_str!(concat!(env!("CARGO_MANIFEST_DIR"), "/templates/home.html")),
            ),
            (
                "chat.html",
                include_str!(concat!(env!("CARGO_MANIFEST_DIR"), "/templates/chat.html")),
            ),
            (
                "login.html",
                include_str!(concat!(env!("CARGO_MANIFEST_DIR"), "/templates/login.html")),
            ),
            (
                "signup.html",
                include_str!(concat!(
                    env!("CARGO_MANIFEST_DIR"),
                    "/templates/signup.html"
                )),
            ),
            // (
            //     "livestream.html",
            //     include_str!(concat!(
            //         env!("CARGO_MANIFEST_DIR"),
            //         "/templates/livestream.html"
            //     )),
            // ),
        ])
        .unwrap();

        let totally_secure_code = b"abcdefghijklmnopqrstuvwxyz123456789";

        App::new()
            .wrap(IdentityService::new(
                CookieIdentityPolicy::new(totally_secure_code) // TODO: Generate actual key
                    .name("session")
                    .max_age(2_629_800) // one month
                    .secure(false),
            ))
            .wrap(middleware::Logger::new("%a \"%r\" %s %b \"%{Referer}i\""))
            .wrap(middleware::NormalizePath::default()) // this doesn't really work properly what the fuck
            .service(actix_files::Files::new(
                "/static",
                concat!(env!("CARGO_MANIFEST_DIR"), "/acolyte-web/dist/"),
            ))
            .data(tera)
            .data(pool.clone())
            .data(srv.clone())
            .service(frontpage::index)
            .service(auth::login)
            .service(auth::login_form)
            .service(auth::signup)
            .service(auth::signup_form)
            .service(
                web::scope("/chat")
                    .service(chat::frontend)
                    .service(chat::ws_upgrader),
            )
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}
