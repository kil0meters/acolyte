-- Your SQL goes here
CREATE TYPE permission_level AS ENUM (
  'AUTH_ADMIN',
  'AUTH_MODERATOR',
  'AUTH_STANDARD',
  'AUTH_BANNED'
)
