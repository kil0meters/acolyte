use std::time;

use actix::Addr;
use actix_web::{error, get, web, Error, HttpRequest, HttpResponse};
use actix_web_actors::ws;

pub mod message_types;
pub mod server;
pub mod session;

#[get("/ws")]
pub async fn ws_upgrader(
    request: HttpRequest,
    stream: web::Payload,
    srv: web::Data<Addr<server::Server>>,
) -> Result<HttpResponse, Error> {
    println!("Making connection to client...");
    let res = ws::start(
        session::Client {
            id: 0,
            username: Some("kilometers".to_owned()),
            auth_level: message_types::AuthLevel::Standard,
            hb: time::Instant::now(),
            conn: srv.get_ref().clone(),
        },
        &request,
        stream,
    );
    println!("{:?}", res);
    res
}

#[get("/")]
pub async fn frontend(
    request: HttpRequest,
    tmpl: web::Data<tera::Tera>,
) -> Result<HttpResponse, Error> {
    println!("Heyo");
    let mut ctx = tera::Context::new();

    ctx.insert("title", "Chat");
    ctx.insert("username", "kilometers");
    ctx.insert("is_embed", &false);

    let s = tmpl
        .render("chat.html", &ctx)
        .map_err(|_| error::ErrorInternalServerError("Template error"))?;

    Ok(HttpResponse::Ok().content_type("text/html").body(s))
}
