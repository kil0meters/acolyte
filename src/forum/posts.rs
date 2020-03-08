use actix_web::web::Form;
use anyhow::Result;
use diesel::pg::PgConnection;
use diesel::prelude::*;

use crate::forum::PostForm;
use crate::models::{Account, Post};
use crate::schema::posts;

const PAGE_LENGTH: i64 = 25;

pub fn get_hot_posts(page_number: i64, conn: &PgConnection) -> Result<Vec<Post>> {
    let recent_posts = posts::table
        .select(posts::all_columns)
        .order_by(posts::created_at)
        .offset(page_number * PAGE_LENGTH)
        .limit(PAGE_LENGTH)
        .load::<Post>(conn)?;

    println!("{:?}", recent_posts);

    Ok(recent_posts)
}

pub fn create_new_post(
    by: Account,
    title: String,
    body: String,
    link: String,
    conn: &PgConnection,
) -> Result<Post> {
    let post = Post::new(by, title, body, link);

    post.validate()?;

    Ok(diesel::insert_into(posts::table)
        .values(&post)
        .get_result::<Post>(conn)?)
}
