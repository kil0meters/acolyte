CREATE TABLE IF NOT EXISTS users (
  id text UNIQUE PRIMARY KEY,
  username text UNIQUE NOT NULL,
  password_hash text NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
  created_at timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
  permissions integer NOT NULL DEFAULT 3
)
