use anyhow::{anyhow, Result};
use diesel::prelude::*;
use diesel::{sql_function, PgConnection};

use crate::models::{Comment, User};
use crate::schema::comments;
use crate::templates::CommentWidget;

sql_function! {
    fn length(x: diesel::sql_types::Text) -> Integer
}

pub fn get_commnet_tree_from_parent(
    parent_id: &str,
    conn: &PgConnection,
) -> Result<Vec<CommentWidget>> {
    let comments = comments::table
        .select(comments::all_columns)
        .filter(comments::id_parents.like(parent_id.to_owned() + "%"))
        .order_by(length(comments::id))
        .load::<Comment>(conn)?;

    let first = match comments.get(0) {
        Some(x) => x,
        None => return Err(anyhow!("No comments")),
    };

    let comment_tree = CommentWidget {
        comment: first.to_owned(),
        has_more_children: false,
        children: vec![],
    };

    Ok(vec![comment_tree])
}

pub fn create_new_comment(
    author: User,
    parent_id: String,
    body: String,
    conn: &PgConnection,
) -> Result<Comment> {
    let comment = Comment::new(author, parent_id, body);

    comment.validate()?;
    Ok(diesel::insert_into(comments::table)
        .values(&comment)
        .get_result::<Comment>(conn)?)
}
