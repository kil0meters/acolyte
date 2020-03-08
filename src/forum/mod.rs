use actix_identity::Identity;
use actix_web::{get, http, post, web, HttpResponse};
use anyhow::anyhow;
// use anyhow::{anyhow, Result};
use askama::Template;
use serde::Deserialize;

use crate::{templates, DbPool};

pub mod posts;

#[derive(Deserialize)]
struct PostForm {
    title: String,
    link: String,
    body: String,
}

#[get("")]
pub async fn index(
    pool: web::Data<DbPool>,
    id: Identity,
) -> Result<HttpResponse, actix_web::Error> {
    let conn = pool.get().expect("Error getting database");

    let posts = web::block(move || posts::get_hot_posts(0, &conn))
        .await
        .map_err(|_| HttpResponse::InternalServerError())?;

    let s = templates::ForumFrontpage {
        posts,
        logged_in: id.identity().is_some(),
    }
    .render()
    .map_err(|_| HttpResponse::InternalServerError())?;

    Ok(HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(s))
}

#[post("/create-post")]
pub async fn submit_post(
    id: Identity,
    pool: web::Data<DbPool>,
    form: web::Form<PostForm>,
) -> Result<HttpResponse, actix_web::Error> {
    if let Some(id) = id.identity() {
        let conn = pool.get().expect("Error getting database");

        let account = serde_json::from_str(&id).unwrap();

        let post = web::block(move || {
            posts::create_new_post(
                account,
                form.title.to_owned(),
                form.body.to_owned(),
                form.link.to_owned(),
                &conn,
            )
        })
        .await
        .map_err(|_| HttpResponse::InternalServerError())?;

        Ok(HttpResponse::SeeOther()
            .header(http::header::LOCATION, format!("/forum/{}", post.id))
            .finish())
    } else {
        Ok(HttpResponse::Unauthorized().finish())
    }
}

#[get("/create-post")]
pub async fn post_editor(id: Identity) -> HttpResponse {
    let s = templates::PostEditor {
        logged_in: id.identity().is_some(),
    }
    .render()
    .unwrap();

    HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(s)
}
