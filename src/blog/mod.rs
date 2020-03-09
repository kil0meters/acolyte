use std::str::FromStr;

use actix_identity::Identity;
use actix_web::{get, http, post, web, Error, HttpResponse};
use anyhow::anyhow;
use askama::Template;
use diesel::prelude::*;
use pulldown_cmark::{html, Options, Parser};
use serde::Deserialize;

use crate::auth::permissions;
use crate::frontpage::HEADER_LINKS;
use crate::models;
use crate::schema::blog_posts;
use crate::templates;
use crate::DbPool;

#[get("/blog")]
pub async fn blog_index(pool: web::Data<DbPool>) -> Result<HttpResponse, Error> {
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
    if let Some(id) = id.identity() {
        let account = serde_json::from_str::<models::Account>(&id).unwrap();

        if account.permissions == permissions::ADMIN {
            match models::BlogPost::from_str(&form.data) {
                Ok(post) => {
                    let conn = pool.get().unwrap();
                    let post_id = post.id.clone();

                    web::block(move || {
                        diesel::insert_into(blog_posts::table)
                            .values(&post)
                            .execute(&conn)
                    })
                    .await
                    .unwrap();

                    HttpResponse::SeeOther()
                        .header(http::header::LOCATION, format!("/blog/{}", post_id))
                        .finish()
                }
                Err(e) => {
                    error!("Encountered an error while parsing blog post: {}", e);

                    HttpResponse::SeeOther()
                        .header(http::header::LOCATION, "/blog/new?error=1")
                        .finish()
                }
            }
        } else {
            HttpResponse::SeeOther()
                .header(http::header::LOCATION, "/unauthorized")
                .finish()
        }
    } else {
        HttpResponse::SeeOther()
            .header(http::header::LOCATION, "/unauthorized")
            .finish()
    }
}

#[get("/blog/new")]
pub async fn blog_editor(id: Identity) -> HttpResponse {
    if let Some(id) = id.identity() {
        let account = serde_json::from_str::<models::Account>(&id).unwrap();

        if account.permissions == permissions::ADMIN {
            let s = templates::BlogPostEditor {}.render().unwrap();

            return HttpResponse::Ok()
                .content_type("text/html; charset=utf-8")
                .body(s);
        }
    }
    HttpResponse::SeeOther()
        .header(http::header::LOCATION, "/unauthorized")
        .finish()
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
