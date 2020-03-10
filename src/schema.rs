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
    threads (id) {
        id -> Text,
        user_id -> Text,
        username -> Text,
        title -> Text,
        body -> Nullable<Text>,
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

joinable!(threads -> users (user_id));

allow_tables_to_appear_in_same_query!(
    blog_posts,
    threads,
    users,
);
