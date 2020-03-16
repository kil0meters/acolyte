use actix_web::{get, HttpResponse};
use askama::Template;

use crate::serve_template;
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
async fn index() -> HttpResponse {
    serve_template!(Homepage {
        live_status: false,
        header_links: &HEADER_LINKS,
    })
}
