use askama::Template;
use serde::{Deserialize, Serialize};

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
pub struct ChatEmbed<'a> {
    pub username: &'a str,
    pub is_embed: bool,
    pub is_moderator: bool,
}
