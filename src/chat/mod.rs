use std::time;

use actix::Addr;
use actix_identity::Identity;
use actix_web::{get, web, Error, HttpRequest, HttpResponse};
use actix_web_actors::ws;
use askama::Template;

use crate::models::User;
use crate::serve_template;
use crate::templates;

pub mod commands;
pub mod message_types;
pub mod server;
pub mod session;

#[get("/ws")]
pub async fn ws_upgrader(
    id: Identity,
    request: HttpRequest,
    stream: web::Payload,
    srv: web::Data<Addr<server::Server>>,
) -> Result<HttpResponse, Error> {
    let user = User::from_identity(id);

    ws::start(
        session::Client {
            id: 0,
            username: user.username,
            auth_level: user.permissions,
            hb: time::Instant::now(),
            conn: srv.get_ref().clone(),
        },
        &request,
        stream,
    )
}

#[get("")]
pub async fn frontend(id: Identity) -> HttpResponse {
    serve_template!(templates::ChatPage {
        user: User::from_identity(id),
        is_embed: false,
    });
}
