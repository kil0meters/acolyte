use std::time;

use actix::Addr;
use actix_identity::Identity;
use actix_web::{error, get, web, Error, HttpRequest, HttpResponse};
use actix_web_actors::ws;
use serde_json::json;

use crate::auth::permissions;
use crate::models::Account;

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
    let username = if let Some(id_data) = id.identity() {
        match serde_json::from_str::<Account>(&id_data) {
            Ok(account) => Some(account.username),
            _ => {
                id.forget();
                None
            }
        }
    } else {
        None
    };

    let res = ws::start(
        session::Client {
            id: 0,
            username,
            auth_level: permissions::LOGGED_OUT,
            hb: time::Instant::now(),
            conn: srv.get_ref().clone(),
        },
        &request,
        stream,
    );
    println!("{:?}", res);
    res
}

#[get("")]
pub async fn frontend(
    // request: HttpRequest,
    id: Identity,
    tmpl: web::Data<tera::Tera>,
) -> Result<HttpResponse, Error> {
    let username = if let Some(id) = id.identity() {
        let account: Account = serde_json::from_str(&id).unwrap();
        account.username
    } else {
        "ANON".to_owned()
    };

    let ctx = tera::Context::from_value(json!({
        "title": "Chat",
        "username": username,
        "is_embed": false
    }))
    .unwrap();

    let s = tmpl
        .render("chat.html", &ctx)
        .map_err(|_| error::ErrorInternalServerError("Template error"))?;

    Ok(HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(s))
}
