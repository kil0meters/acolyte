table! {
    blog_posts (id) {
        id -> Text,
        title -> Text,
        body -> Text,
        updated_at -> Timestamp,
        created_at -> Timestamp,
    }
}

table! {
    comments (id) {
        id -> Text,
        id_parents -> Text,
        user_id -> Text,
        username -> Text,
        body -> Text,
        body_html -> Text,
        removed -> Bool,
        updated_at -> Timestamp,
        created_at -> Timestamp,
        upvotes -> Int4,
        downvotes -> Int4,
    }
}

table! {
    threads (id) {
        id -> Text,
        user_id -> Text,
        username -> Text,
        title -> Text,
        body -> Nullable<Text>,
        body_html -> Nullable<Text>,
        link -> Nullable<Text>,
        removed -> Bool,
        updated_at -> Timestamp,
        created_at -> Timestamp,
        upvotes -> Int4,
        downvotes -> Int4,
    }
}

table! {
    users (id) {
        id -> Text,
        username -> Text,
        password_hash -> Text,
        updated_at -> Timestamp,
        created_at -> Timestamp,
        permissions -> Int4,
    }
}

allow_tables_to_appear_in_same_query!(
    blog_posts,
    comments,
    threads,
    users,
);
