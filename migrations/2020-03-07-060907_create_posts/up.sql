CREATE TABLE IF NOT EXISTS posts (
  id text UNIQUE PRIMARY KEY,
  account_id text NOT NULL REFERENCES accounts(id),
  title text NOT NULL,
  body text,
  link text,
  removed boolean NOT NULL DEFAULT FALSE,
  created_at timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
  upvotes integer NOT NULL DEFAULT 0,
  downvotes integer NOT NULL DEFAULT 0
)
