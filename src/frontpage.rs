use actix_web::{error, get, web, Error, HttpResponse};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
struct HeaderLink<'a> {
    title: &'a str,
    url: &'a str,
}

impl<'a> HeaderLink<'a> {
    fn homepage() -> Vec<HeaderLink<'a>> {
        vec![
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
        ]
    }
}

#[get("/")]
async fn index(tmpl: web::Data<tera::Tera>) -> Result<HttpResponse, Error> {
    let mut ctx = tera::Context::new();

    ctx.insert("title", "Miles Benton");
    ctx.insert("page_title", "Miles Benton");
    ctx.insert("live_status", &true);
    ctx.insert("header", &HeaderLink::homepage());

    let s = tmpl
        .render("home.html", &ctx)
        .map_err(|_| error::ErrorInternalServerError("Template error"))?;

    Ok(HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(s))
}
