use actix_identity::Identity;
use actix_web::{get, post, web, Error, HttpResponse};
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
pub async fn index(pool: web::Data<DbPool>, id: Identity) -> Result<HttpResponse, Error> {
    let conn = pool.get().expect("Error getting database");

    let posts = web::block(move || posts::get_hot_posts(0, &conn))
        .await
        .unwrap();

    let s = templates::ForumFrontpage {
        posts,
        logged_in: id.identity().is_some(),
    }
    .render()
    .unwrap();

    Ok(HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(s))
}

// #[post("/create-post")]
// pub async fn submit_post(
//     id: Identity,
//     pool: web::Data<DbPool>,
//     form: web::Form<PostForm>,
// ) -> HttpResponse {
// }

// #[get("/create-post")]
// pub async fn post_editor() -> HttpResponse {}
