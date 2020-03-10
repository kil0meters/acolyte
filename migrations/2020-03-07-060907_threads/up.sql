CREATE TABLE IF NOT EXISTS threads (
  id text UNIQUE PRIMARY KEY,
  user_id text NOT NULL REFERENCES users(id),
  username text NOT NULL,
  title text NOT NULL,
  body text,
  link text,
  removed boolean NOT NULL DEFAULT FALSE,
  updated_at timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
  created_at timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
  upvotes integer NOT NULL DEFAULT 0,
  downvotes integer NOT NULL DEFAULT 0
)
