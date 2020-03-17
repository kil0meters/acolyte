#[macro_export]
macro_rules! not_found {
    () => {{
        use askama::Template;
        $crate::serve_template!($crate::templates::NotFound { header_links: &[] })
    }};
}

#[macro_export]
macro_rules! unauthorized {
    () => {{
        use askama::Template;
        $crate::serve_template!($crate::templates::Unauthorized { header_links: &[] })
    }};
}

#[macro_export]
macro_rules! redirect_to {
    ($e:expr) => {
        return HttpResponse::SeeOther()
            .header(http::header::LOCATION, $e)
            .finish();
    };

    ($($arg:tt)+) => {
        return HttpResponse::SeeOther()
            .header(http::header::LOCATION, format!($($arg)+))
            .finish();
    };
}

#[macro_export]
macro_rules! unwrap_or_notfound {
    ($e:expr) => {
        match $e {
            Ok(x) => x,
            Err(_) => $crate::not_found!(),
        }
    };
}

#[macro_export]
macro_rules! unwrap_or_redirect {
    ($e:expr => $target:expr) => {
        match $e {
            Ok(x) => x,
            Err(_) => $crate::redirect_to!($target),
        }
    };
}

/// Serves html from a string
#[macro_export]
macro_rules! serve_html {
    ($e:expr) => {
        return HttpResponse::Ok()
            .content_type("text/html; charset=utf-8")
            .body($e);
    };
}

#[macro_export]
macro_rules! serve_template {
    ($e:expr) => {
        $crate::serve_html!($e.render().unwrap());
    };
}
