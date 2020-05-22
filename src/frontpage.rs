use actix_identity::Identity;
use actix_web::{get, HttpResponse};
use askama::Template;

use crate::models::User;
use crate::serve_template;
use crate::templates::{HeaderLink, Homepage, StreamEmbed};

pub const HEADER_LINKS: [HeaderLink<'static>; 6] = [
    HeaderLink {
        title: "Forum",
        url: "/forum",
    },
    HeaderLink {
        title: "Videos",
        url: "/forum",
    },
    HeaderLink {
        title: "Live",
        url: "/live",
    },
    HeaderLink {
        title: "Logs",
        url: "/logs",
    },
    HeaderLink {
        title: "Blog",
        url: "/blog",
    },
    HeaderLink {
        title: "Resume",
        url: "/resume.pdf",
    },
];

pub const CHANNEL_ID: &'static str = "UCSJ4gkVC6NrvII8umztf0Ow";

#[get("/")]
async fn index() -> HttpResponse {
    serve_template!(Homepage {
        live_status: false,
        header_links: &HEADER_LINKS,
    })
}

#[get("/live")]
async fn stream_embed(id: Identity) -> HttpResponse {
    let user = User::from_identity(id);

    serve_template!(StreamEmbed {
        user,
        channel_id: CHANNEL_ID,
    })
}
