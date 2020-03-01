table! {
    accounts (id) {
        id -> Text,
        username -> Text,
        password_hash -> Text,
        created_at -> Timestamp,
        permissions -> Int4,
    }
}
