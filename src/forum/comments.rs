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

    let tree = build_tree(comments);
    Ok(tree)
}

// this is probably poorly written, but is still significantly better than
// than doing recursive sql loops
fn build_tree(comments: Vec<Comment>) -> Vec<CommentWidget> {
    // end recursion
    if comments.len() == 0 {
        return Vec::new();
    }

    let mut tree = Vec::new();
    let mut current_shortest_length = usize::max_value();
    for comment in comments {
        let id_len = comment.id_parents.len();

        if id_len <= current_shortest_length {
            current_shortest_length = id_len;

            tree.push(CommentWidget {
                comment,
                has_more_children: true,
                children: Vec::new(),
            });
        } else {
            for parent_element in &mut tree {
                if comment
                    .id_parents
                    .starts_with(&parent_element.comment.id_parents)
                {
                    parent_element.children.push(CommentWidget {
                        comment,
                        has_more_children: true,
                        children: Vec::new(),
                    });

                    break;
                }
            }
        }
    }

    for mut element in &mut tree {
        element.children = build_tree(
            element
                .children
                .iter()
                .map(|x| x.comment.clone())
                .collect::<Vec<Comment>>(),
        );

        if element.children.len() == 0 {
            element.has_more_children = false;
        }
    }

    tree
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
