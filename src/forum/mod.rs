use actix_identity::Identity;
use actix_web::{get, http, post, web, HttpResponse};
use diesel::prelude::*;
// use anyhow::{anyhow, Result};
use askama::Template;
use serde::Deserialize;

use crate::auth::permissions;
use crate::auth::permissions::Permission;
use crate::models::{Thread, User};
use crate::serve_template;
use crate::{redirect_to, templates, unauthorized, unwrap_or_notfound, unwrap_or_redirect, DbPool};

mod comments;
pub mod threads;

#[derive(Deserialize)]
pub struct ThreadForm {
    title: String,
    link: String,
    body: String,
}

#[get("")]
pub async fn index(pool: web::Data<DbPool>, id: Identity) -> HttpResponse {
    let conn = pool.get().expect("Error getting database");

    let threads = web::block(move || threads::get_hot_threads(0, &conn))
        .await
        .unwrap_or(vec![]);

    serve_template!(templates::ForumFrontpage {
        user: User::from_identity(id),
        threads,
    });
}

#[post("/create-thread")]
pub async fn create_thread_form(
    id: Identity,
    pool: web::Data<DbPool>,
    form: web::Form<ThreadForm>,
) -> HttpResponse {
    let user = User::from_identity(id);

    if !user.permissions.at_least(permissions::STANDARD) {
        unauthorized!();
    }

    let conn = pool.get().expect("Error getting database");

    let thread = unwrap_or_redirect!({
        web::block(move || {
            threads::create_new_thread(
                user,
                form.title.to_owned(),
                form.body.to_owned(),
                form.link.to_owned(),
                &conn,
            )
        })
        .await
    } => "/create-thread?error=1");

    redirect_to!("/forum/thread/{}", thread.id);
}

#[get("/create-thread")]
pub async fn thread_editor(id: Identity) -> HttpResponse {
    let user = User::from_identity(id);

    serve_template!(templates::ThreadEditor { user });
}

#[derive(Deserialize, Debug)]
pub struct CommentForm {
    thread_id: String,
    parent_id: String,
    body: String,
}

#[post("/create-comment")]
pub async fn comment_form(
    id: Identity,
    form: web::Form<CommentForm>,
    pool: web::Data<DbPool>,
) -> HttpResponse {
    debug!("New comment: {:?}", form);
    let user = User::from_identity(id);

    if !user.permissions.at_least(permissions::STANDARD) {
        unauthorized!();
    }

    let conn = pool.get().unwrap();

    let thread_id = form.thread_id.to_owned();
    let comment = unwrap_or_notfound!({
        web::block(move || {
            comments::create_new_comment(
                user,
                form.parent_id.to_owned(),
                form.body.to_owned(),
                &conn,
            )
        })
        .await
    });

    redirect_to!("/forum/thread/{}#{}", thread_id, comment.id);
}

#[get("/thread/{id}")]
pub async fn serve_thread(
    id: Identity,
    pool: web::Data<DbPool>,
    thread_id: web::Path<String>,
) -> HttpResponse {
    let user = User::from_identity(id);

    let conn = pool.get().unwrap();
    let (thread, comments) = unwrap_or_notfound!(
        web::block(move || {
            use crate::schema::threads;

            let thread = threads::table
                .select(threads::all_columns)
                .filter(threads::id.eq(&thread_id.to_string()))
                .first::<Thread>(&conn)?;

            let comments = comments::get_commnet_tree_from_parent(&thread_id.to_string(), &conn)
                .unwrap_or(vec![]);

            Ok::<(Thread, Vec<templates::CommentWidget>), anyhow::Error>((thread, comments))
        })
        .await
    );

    serve_template!(templates::ForumThread {
        user,
        thread,
        comments,
    })
}

