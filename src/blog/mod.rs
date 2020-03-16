use std::str::FromStr;

use actix_identity::Identity;
use actix_web::{get, http, post, web, HttpResponse};
use anyhow::anyhow;
use askama::Template;
use diesel::prelude::*;
use pulldown_cmark::{html, Options, Parser};
use serde::Deserialize;

use crate::auth::permissions;
use crate::frontpage::HEADER_LINKS;
use crate::models::{BlogPost, User};
use crate::schema::blog_posts;
use crate::templates::{self, BlogIndex, BlogPostEditor};
use crate::DbPool;
use crate::{redirect_to, serve_template, unauthorized, unwrap_or_notfound, unwrap_or_redirect};

#[get("/blog")]
pub async fn blog_index(pool: web::Data<DbPool>) -> HttpResponse {
    let conn = pool.get().unwrap();

    let posts = web::block(move || {
        blog_posts::table
            .select(blog_posts::all_columns)
            .load::<BlogPost>(&conn)
            .map_err(|_| anyhow!("Error fetching blogpost"))
    })
    .await
    .unwrap_or(vec![]);

    serve_template!(BlogIndex {
        posts,
        header_links: &HEADER_LINKS,
    });
}

#[derive(Deserialize)]
pub struct BlogUploadForm {
    data: String,
}

// TODO: create `redirect_if!` macro
#[post("/blog/new")]
pub async fn blog_upload_form(
    id: Identity,
    pool: web::Data<DbPool>,
    form: web::Form<BlogUploadForm>,
) -> HttpResponse {
    let user = User::from_identity(id);

    if user.permissions != permissions::ADMIN {
        unauthorized!();
    }

    let post = unwrap_or_redirect!(
        BlogPost::from_str(&form.data) => "/blog/new?error=1");

    let conn = pool.get().unwrap();
    let post_id = post.id.clone();

    unwrap_or_redirect!({
        web::block(move || {
            diesel::insert_into(blog_posts::table)
                .values(&post)
                .execute(&conn)
        }).await
    } => "/blog/new?error=1");

    redirect_to!("/blog/{}", post_id);
}

#[get("/blog/new")]
pub async fn blog_editor(id: Identity) -> HttpResponse {
    let user = User::from_identity(id);

    if user.permissions != permissions::ADMIN {
        unauthorized!();
    }

    serve_template!(BlogPostEditor {});
}

#[get("/blog/{id}")]
pub async fn serve_blog_post(pool: web::Data<DbPool>, id: web::Path<String>) -> HttpResponse {
    let id = id.to_string();
    let conn = pool.get().unwrap();

    let post: BlogPost = unwrap_or_notfound!(
        web::block(move || {
            blog_posts::table
                .select(blog_posts::all_columns)
                .filter(blog_posts::id.eq(&id))
                .first::<BlogPost>(&conn)
        })
        .await
    );

    // let (content, metadata) = extract_metadata(post.body).unwrap();
    let html_output = render_content(&post.body);

    serve_template!(templates::BlogPost {
        title: post.title,
        created_at: post.created_at,
        updated_at: post.updated_at,
        post_body: html_output,
        header_links: &HEADER_LINKS,
    });
}

fn render_content(content: &str) -> String {
    let parser = Parser::new_ext(content, Options::all());

    let mut html_output = String::new();
    html::push_html(&mut html_output, parser);

    html_output
}
