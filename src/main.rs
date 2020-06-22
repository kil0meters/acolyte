#![warn(clippy::all)]

#[macro_use]
extern crate diesel;

#[macro_use]
extern crate lazy_static;

#[macro_use]
extern crate log;

use std::borrow::Cow;
use std::env;

use actix::Actor;
use actix_identity::{CookieIdentityPolicy, IdentityService};
use actix_web::body::Body;
use actix_web::{http, middleware, web, App, HttpRequest, HttpResponse, HttpServer};
use chrono::Utc;
use clap::Clap;
use diesel::pg::PgConnection;
use diesel::r2d2::{ConnectionManager, Pool};
use dotenv::dotenv;
use rust_embed::RustEmbed;

pub mod auth;
pub mod blog;
pub mod chat;
pub mod forum;
pub mod frontpage;
pub mod macros;
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
    // #[clap(short, long, default_value = "https://localhost:8000")]
    // /// Not implemented yet. Would be used for pooling chat instances together.
    // master: String,
}

#[derive(Clap)]
/// A standalone forum instance.
struct Forum {
    #[clap(short, long, default_value = "8080")]
    /// Port to listen on
    port: String,
    // #[clap(short, long, default_value = "http://localhost:8080")]
    // /// Link to chat instance
    // chat_url: String,
}

#[derive(RustEmbed)]
#[folder = "acolyte-web/dist/"]
struct Asset;

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

// TODO: this function is poorly written; duplicates a lot of code
fn handle_embedded_file(path: &str, accepted_encodings: &str) -> HttpResponse {
    let path = path.to_owned();

    if accepted_encodings.contains("br") {
        if let Some(content) = Asset::get(&(path.clone() + ".br")) {
            let body: Body = match content {
                Cow::Borrowed(bytes) => bytes.into(),
                Cow::Owned(bytes) => bytes.into(),
            };

            return HttpResponse::Ok()
                .content_type(mime_guess::from_path(path).first_or_octet_stream().as_ref())
                .header(http::header::CONTENT_ENCODING, "br")
                .body(body);
        }
    }

    if accepted_encodings.contains("gzip") {
        if let Some(content) = Asset::get(&(path.clone() + ".gz")) {
            let body: Body = match content {
                Cow::Borrowed(bytes) => bytes.into(),
                Cow::Owned(bytes) => bytes.into(),
            };

            return HttpResponse::Ok()
                .content_type(mime_guess::from_path(path).first_or_octet_stream().as_ref())
                .header(http::header::CONTENT_ENCODING, "gzip")
                .body(body);
        }
    }

    match Asset::get(&path) {
        Some(content) => {
            let body: Body = match content {
                Cow::Borrowed(bytes) => bytes.into(),
                Cow::Owned(bytes) => bytes.into(),
            };

            HttpResponse::Ok()
                .content_type(mime_guess::from_path(path).first_or_octet_stream().as_ref())
                .body(body)
        }
        None => HttpResponse::NotFound().body("404 Not Found"),
    }
}

fn dist(req: HttpRequest) -> HttpResponse {
    let encodings = req
        .headers()
        .get(http::header::ACCEPT_ENCODING)
        .unwrap()
        .to_str()
        .unwrap_or("");

    let path = &req.path()["/static/".len()..]; // trim the preceding `/static/` in path
    handle_embedded_file(path, encodings)
}

#[actix_rt::main]
async fn main() -> std::io::Result<()> {
    // parse arguments
    let opts = Opts::parse();

    // std::env::set_var("RUST_LOG", "actix_web=info");
    setup_logger(opts.log_level).unwrap();

    dotenv().ok();
    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL is required.");
    let manager = ConnectionManager::<PgConnection>::new(database_url);
    let pool = Pool::builder()
        .build(manager)
        .expect("Failed to create database pool");

    let logger_format = "%a \"%r\" %s %b \"%{Referer}i\" %Dms".to_owned();
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
                    .wrap(middleware::Logger::new(&logger_format))
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
                    .wrap(middleware::Logger::new(&logger_format))
                    .wrap(middleware::NormalizePath::default()) // this doesn't really work properly what the fuck
                    .data(pool.clone())
                    .service(web::resource("/static/{_:.*}").route(web::get().to(dist)))
                    .service(frontpage::index)
                    .service(auth::login)
                    .service(auth::login_form)
                    .service(auth::signup)
                    .service(auth::signup_form)
                    .service(
                        web::scope("/forum")
                            .service(forum::index)
                            .service(forum::create_thread_form)
                            .service(forum::thread_editor),
                    )
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
                    .wrap(middleware::Logger::new(&logger_format))
                    .wrap(middleware::NormalizePath::default()) // this doesn't really work properly what the fuck
                    .service(web::resource("/static/{_:.*}").route(web::get().to(dist)))
                    .data(pool.clone())
                    .data(srv.clone())
                    .service(frontpage::index)
                    .service(frontpage::stream_embed)
                    .service(auth::login)
                    .service(auth::login_form)
                    .service(auth::signup)
                    .service(auth::signup_form)
                    .service(blog::blog_index)
                    .service(blog::blog_editor)
                    .service(blog::blog_upload_form)
                    .service(blog::serve_blog_post)
                    .service(
                        web::scope("/forum")
                            .service(forum::index)
                            .service(forum::create_thread_form)
                            .service(forum::thread_editor)
                            .service(forum::comment_form)
                            .service(forum::serve_thread),
                    )
                    .service(
                        web::scope("/chat")
                            .service(chat::frontend)
                            .service(chat::ws_upgrader),
                    )
                    .default_service(web::get().to(|| not_found!()))
            })
            .bind(format!("0.0.0.0:{}", opts.port))?
            .run()
            .await;
        }
    };
}
