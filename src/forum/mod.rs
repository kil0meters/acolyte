use actix_identity::Identity;
use actix_web::{get, http, post, web, Error, HttpResponse};
use diesel::prelude::*;
// use anyhow::{anyhow, Result};
use askama::Template;
use serde::Deserialize;

use crate::auth::permissions;
use crate::auth::permissions::Permission;
use crate::models::{Thread, User};
use crate::{templates, DbPool};

pub mod threads;

#[derive(Deserialize)]
pub struct ThreadForm {
    title: String,
    link: String,
    body: String,
}

#[get("")]
pub async fn index(pool: web::Data<DbPool>, id: Identity) -> Result<HttpResponse, Error> {
    let conn = pool.get().expect("Error getting database");

    let threads = web::block(move || threads::get_hot_threads(0, &conn))
        .await
        .map_err(|_| HttpResponse::InternalServerError())?;

    let s = templates::ForumFrontpage {
        user: User::from_identity(id),
        threads,
    }
    .render()
    .map_err(|_| HttpResponse::InternalServerError())?;

    Ok(HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(s))
}

#[post("/create-thread")]
pub async fn create_thread_form(
    id: Identity,
    pool: web::Data<DbPool>,
    form: web::Form<ThreadForm>,
) -> Result<HttpResponse, Error> {
    let user = User::from_identity(id);

    if !user.permissions.at_least(permissions::STANDARD) {
        debug!("Unauthorized request");
        return Ok(HttpResponse::Unauthorized().finish());
    }

    let conn = pool.get().expect("Error getting database");

    let thread = web::block(move || {
        threads::create_new_thread(
            user,
            form.title.to_owned(),
            form.body.to_owned(),
            form.link.to_owned(),
            &conn,
        )
    })
    .await
    .map_err(|e| {
        error!("error creating thread: {}", e);
        HttpResponse::InternalServerError();
    })?;

    Ok(HttpResponse::SeeOther()
        .header(http::header::LOCATION, format!("/forum/{}", thread.id))
        .finish())
}

#[get("/create-thread")]
pub async fn thread_editor(id: Identity) -> HttpResponse {
    let s = templates::ThreadEditor {
        user: User::from_identity(id),
    }
    .render()
    .unwrap();

    HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(s)
}

#[get("/{id}")]
pub async fn serve_thread(
    id: Identity,
    pool: web::Data<DbPool>,
    post_id: web::Path<String>,
) -> Result<HttpResponse, Error> {
    use crate::schema::threads;

    let user = User::from_identity(id);

    let conn = pool
        .get()
        .map_err(|_| HttpResponse::InternalServerError())?;

    let thread = web::block(move || {
        threads::table
            .select(threads::all_columns)
            .filter(threads::id.eq(&post_id.to_string()))
            .first::<Thread>(&conn)
    })
    .await
    .map_err(|_| HttpResponse::InternalServerError())?;

    let s = templates::ForumThread { user, thread }.render().unwrap();

    Ok(HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(s))
}
