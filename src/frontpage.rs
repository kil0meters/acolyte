use actix_web::{error, get, web, Error, HttpResponse};
use askama::Template;

use crate::templates::{HeaderLink, Homepage};

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

#[get("/")]
async fn index() -> Result<HttpResponse, Error> {
    let t = Homepage {
        live_status: false,
        header_links: &HEADER_LINKS,
    }
    .render()
    .unwrap();

    Ok(HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(t))
}
