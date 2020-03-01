CREATE TABLE IF NOT EXISTS accounts (
  id text UNIQUE PRIMARY KEY,
  username text UNIQUE NOT NULL,
  password_hash text NOT NULL,
  created_at timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
  permissions integer NOT NULL DEFAULT 3
)
