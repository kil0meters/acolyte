table! {
    accounts (id) {
        id -> Text,
        username -> Text,
        password_hash -> Text,
        created_at -> Timestamp,
        permissions -> Int4,
    }
}

table! {
    posts (id) {
        id -> Text,
        account_id -> Text,
        title -> Text,
        body -> Nullable<Text>,
        link -> Nullable<Text>,
        removed -> Bool,
        created_at -> Timestamp,
        upvotes -> Int4,
        downvotes -> Int4,
    }
}

joinable!(posts -> accounts (account_id));

allow_tables_to_appear_in_same_query!(
    accounts,
    posts,
);
