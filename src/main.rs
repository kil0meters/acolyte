#![warn(clippy::all)]

#[macro_use]
extern crate diesel;

#[macro_use]
extern crate lazy_static;

use log;
use std::{env, time};

use actix::Actor;
use actix_identity::{CookieIdentityPolicy, IdentityService};
use actix_web::{middleware, web, App, HttpServer};
use chrono::Utc;
use clap::Clap;
use diesel::pg::PgConnection;
use diesel::r2d2::{ConnectionManager, Pool};

pub mod auth;
pub mod chat;
pub mod frontpage;
pub mod models;
pub mod schema;
pub mod templates;

type DbPool = Pool<ConnectionManager<PgConnection>>;

#[derive(Clap)]
#[clap(version = "0.1.0", author = "kilometers")]
/// A web forum for influencers allowing you to control your audience.
///
/// By default, it will run a bundle of the forum and chat.
struct Opts {
    #[clap(short, long, default_value = "Debug")]
    /// Logging level
    log_level: log::LevelFilter,

    #[clap(short, long, default_value = "8080")]
    /// Port to listen on
    port: String,

    #[clap(subcommand)]
    subcommand: Option<SubCommand>,
}

#[derive(Clap)]
enum SubCommand {
    /// A standalone chat instance.
    #[clap(name = "chat")]
    Chat(Chat),

    /// A standalone chat instance.
    #[clap(name = "forum")]
    Forum(Forum),
}

#[derive(Clap)]
/// A standalone chat instance.
struct Chat {
    #[clap(short, long, default_value = "8080")]
    /// Port to listen on
    port: String,

    #[clap(short, long, default_value = "https://localhost:8000")]
    /// Not implemented yet. Would be used for pooling chat instances together.
    master: String,
}

#[derive(Clap)]
/// A standalone forum instance.
struct Forum {
    #[clap(short, long, default_value = "8080")]
    /// Port to listen on
    port: String,

    #[clap(short, long, default_value = "http://localhost:8080")]
    /// Link to chat instance
    chat_url: String,
}

fn setup_logger(level: log::LevelFilter) -> Result<(), fern::InitError> {
    fern::Dispatch::new()
        .format(|out, message, record| {
            out.finish(format_args!(
                "{}[{}][{}] {}",
                Utc::now().format("[%Y-%m-%d][%H:%M:%S]"),
                record.target(),
                record.level(),
                message
            ))
        })
        .level(level)
        .chain(std::io::stdout())
        .chain(fern::log_file("output.log")?)
        .apply()?;
    Ok(())
}

#[actix_rt::main]
async fn main() -> std::io::Result<()> {
    // parse arguments
    let opts = Opts::parse();

    // std::env::set_var("RUST_LOG", "actix_web=info");
    setup_logger(opts.log_level).unwrap();

    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL is required.");
    let manager = ConnectionManager::<PgConnection>::new(database_url);
    let pool = Pool::builder()
        .build(manager)
        .expect("Failed to create database pool");

    let totally_secure_code = b"abcdefghijklmnopqrstuvwxyz123456789";

    match opts.subcommand {
        // Standalone chat
        Some(SubCommand::Chat(chat)) => {
            let srv = chat::server::Server::default().start();

            return HttpServer::new(move || {
                App::new()
                    .wrap(IdentityService::new(
                        CookieIdentityPolicy::new(totally_secure_code) // TODO: Generate actual key
                            .name("session")
                            .max_age(2_629_800) // one month
                            .secure(false),
                    ))
                    .wrap(middleware::Logger::new("%a \"%r\" %s %b \"%{Referer}i\""))
                    .wrap(middleware::NormalizePath::default()) // this doesn't really work properly what the fuck
                    .data(pool.clone())
                    .data(srv.clone())
                    .service(chat::ws_upgrader)
            })
            .bind(format!("0.0.0.0:{}", chat.port))?
            .run()
            .await;
        }

        // Standalone forum
        Some(SubCommand::Forum(forum)) => {
            return HttpServer::new(move || {
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
                    .data(pool.clone())
                    .service(frontpage::index)
                    .service(auth::login)
                    .service(auth::login_form)
                    .service(auth::signup)
                    .service(auth::signup_form)
            })
            .bind(format!("0.0.0.0:{}", forum.port))?
            .run()
            .await;
        }

        // Default bundled
        None => {
            let srv = chat::server::Server::default().start();

            return HttpServer::new(move || {
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
            .bind(format!("0.0.0.0:{}", opts.port))?
            .run()
            .await;
        }
    };
}
