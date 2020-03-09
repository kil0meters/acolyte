use std::str::FromStr;

use actix_web::{get, web, Error, HttpResponse};
use anyhow::anyhow;
use askama::Template;
use diesel::prelude::*;
use pulldown_cmark::{html, Options, Parser};

use crate::frontpage::HEADER_LINKS;
use crate::models;
use crate::schema::blog_posts;
use crate::templates;
use crate::DbPool;

#[get("/blog")]
pub async fn blog_handler(pool: web::Data<DbPool>) -> Result<HttpResponse, Error> {
    let conn = pool
        .get()
        .map_err(|_| HttpResponse::InternalServerError())?;

    let posts = web::block(move || {
        blog_posts::table
            .select(blog_posts::all_columns)
            .load::<models::BlogPost>(&conn)
            .map_err(|_| anyhow!("Error fetching blogpost"))
    })
    .await
    .map_err(|_| HttpResponse::InternalServerError())?;

    let s = templates::BlogIndex {
        posts,
        header_links: &HEADER_LINKS,
    }
    .render()
    .map_err(|_| HttpResponse::InternalServerError())?;

    Ok(HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(s))
}

#[derive(Form)]
struct BlogUploadForm {
    data: String,
}

#[post("/blog/new")]
pub async fn blog_upload_form(
    id: Identity,
    data: web::Form<BlogUploadForm>,
) -> Result<HttpResponse, Error> {
}

#[get("/blog/new")]
pub async fn blog_upload() -> Result<HttpResponse, Error> {
    Ok(HttpResponse::Ok())
}

#[get("/blog/{id}")]
pub async fn serve_blogpost(
    pool: web::Data<DbPool>,
    id: web::Path<String>,
) -> Result<HttpResponse, Error> {
    let id: String = id.to_string();
    let conn = pool
        .get()
        .map_err(|_| HttpResponse::InternalServerError())?;

    let post = web::block(move || {
        blog_posts::table
            .select(blog_posts::all_columns)
            .filter(blog_posts::id.eq(&id))
            .first::<models::BlogPost>(&conn)
            .map_err(|_| anyhow!("Error fetching blogpost"))
    })
    .await
    .map_err(|_| HttpResponse::InternalServerError())?;

    // let (content, metadata) = extract_metadata(post.body).unwrap();
    let html_output = render_content(&post.body);

    let s = templates::BlogPost {
        title: post.title,
        created_at: post.created_at,
        last_modified: post.last_modified,
        post_body: html_output,
        header_links: &HEADER_LINKS,
    }
    .render()
    .unwrap();

    Ok(HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(s))
}

fn render_content(content: &str) -> String {
    let parser = Parser::new_ext(content, Options::all());

    let mut html_output = String::new();
    html::push_html(&mut html_output, parser);

    html_output
}

/// Regenerates all post markdown, graphs, etc.
fn generate_templates() {}
