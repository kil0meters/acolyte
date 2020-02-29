#![warn(clippy::all)]

use actix::Actor;
use actix_web::{middleware, web, App, HttpServer};
use tera::Tera;

mod chat;
mod frontpage;

#[actix_rt::main]
async fn main() -> std::io::Result<()> {
    println!("Serving on 127.0.0.1:8080");

    std::env::set_var("RUST_LOG", "actix_web=info");
    env_logger::init();

    HttpServer::new(|| {
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
                "livestream.html",
                include_str!(concat!(env!("CARGO_MANIFEST_DIR"), "/templates/home.html")),
            ),
        ])
        .unwrap();

        let srv = chat::server::Server::default().start();

        App::new()
            .wrap(middleware::NormalizePath)
            .wrap(middleware::Logger::default())
            .wrap(
                middleware::DefaultHeaders::new().header("Content-Type", "text/html;charset=utf-8"),
            )
            .service(
                actix_files::Files::new(
                    "/static",
                    concat!(env!("CARGO_MANIFEST_DIR"), "/acolyte-web/dist/"),
                )
                .show_files_listing(),
            )
            .data(tera)
            .data(srv)
            .service(frontpage::index)
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
