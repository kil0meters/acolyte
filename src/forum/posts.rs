use diesel::pg::PgConnection;
use diesel::prelude::*;
use diesel::result::Error;

use crate::models::Post;
use crate::schema::posts;

const PAGE_LENGTH: i64 = 25;

pub fn get_hot_posts(page_number: i64, conn: &PgConnection) -> Result<Vec<Post>, Error> {
    let recent_posts = posts::table
        .select(posts::all_columns)
        .order_by(posts::created_at)
        .offset(page_number * PAGE_LENGTH)
        .limit(PAGE_LENGTH)
        .load::<Post>(conn)?;

    println!("{:?}", recent_posts);

    Ok(recent_posts)
}

// pub fn create_new_post() -> Result<Post, Error> {
// }
