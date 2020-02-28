#![warn(clippy::all)]

use actix_web::{error, get, middleware, web, App, Error, HttpResponse, HttpServer};
use serde::{Deserialize, Serialize};
use tera::Tera;

#[derive(Serialize, Deserialize, Debug)]
struct HeaderLink<'a> {
    title: &'a str,
    url: &'a str,
}

impl<'a> HeaderLink<'a> {
    fn homepage() -> Vec<HeaderLink<'a>> {
        vec![HeaderLink {
            title: "wow",
            url: "https://google.com",
        }]
    }
}

#[get("/")]
async fn index(tmpl: web::Data<tera::Tera>) -> Result<HttpResponse, Error> {
    let mut ctx = tera::Context::new();

    ctx.insert("live_status", &true);
    ctx.insert("header", &HeaderLink::homepage());

    let s = tmpl
        .render("home.html", &ctx)
        .map_err(|_| error::ErrorInternalServerError("Template error"))?;

    Ok(HttpResponse::Ok().content_type("text/html").body(s))
}

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
                "livestream.html",
                include_str!(concat!(env!("CARGO_MANIFEST_DIR"), "/templates/home.html")),
            ),
        ])
        .unwrap();

        App::new()
            .data(tera)
            .wrap(middleware::Logger::default())
            .service(
                actix_files::Files::new(
                    "/static",
                    concat!(env!("CARGO_MANIFEST_DIR"), "/acolyte-web/dist/"),
                )
                .show_files_listing(),
            )
            .service(index)
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}
