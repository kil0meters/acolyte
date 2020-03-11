use askama::Template;
use serde::{Deserialize, Serialize};

// used in templates
use crate::auth::permissions;
use crate::auth::permissions::Permission;

use crate::models;

#[derive(Serialize, Deserialize, Debug, Copy, Clone)]
pub struct HeaderLink<'a> {
    pub title: &'a str,
    pub url: &'a str,
}

#[derive(Template)]
#[template(path = "home.html")]
pub struct Homepage<'a> {
    pub header_links: &'a [HeaderLink<'a>],
    pub live_status: bool,
}

#[derive(Template)]
#[template(path = "signup.html")]
pub struct Signup<'a> {
    pub header_links: &'a [HeaderLink<'a>],
    pub target: &'static str,
    pub error: bool,
}

#[derive(Template)]
#[template(path = "login.html")]
pub struct Login<'a> {
    pub header_links: &'a [HeaderLink<'a>],
    pub target: &'static str,
    pub error: bool,
}

#[derive(Template)]
#[template(path = "chat.html")]
pub struct ChatPage {
    pub user: models::User,
    pub is_embed: bool,
}

#[derive(Template)]
#[template(path = "forum_frontpage.html")]
pub struct ForumFrontpage {
    pub user: models::User,
    pub threads: Vec<models::Thread>,
}

#[derive(Template)]
#[template(path = "thread_editor.html")]
pub struct ThreadEditor {
    pub user: models::User,
}

#[derive(Template)]
#[template(path = "blog_post_page.html", escape = "none")]
pub struct BlogPost<'a> {
    pub title: String,
    pub post_body: String,
    pub created_at: chrono::NaiveDateTime,
    pub updated_at: chrono::NaiveDateTime,
    pub header_links: &'a [HeaderLink<'a>],
}

#[derive(Template)]
#[template(path = "blog_index.html")]
pub struct BlogIndex<'a> {
    pub posts: Vec<models::BlogPost>,
    pub header_links: &'a [HeaderLink<'a>],
}

#[derive(Template)]
#[template(path = "blog_post_editor.html")]
pub struct BlogPostEditor {}

#[derive(Template)]
#[template(path = "forum_thread.html")]
pub struct ForumThread {
    pub user: models::User,
    pub thread: models::Thread,
}
