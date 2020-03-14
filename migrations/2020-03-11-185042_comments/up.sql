CREATE TABLE IF NOT EXISTS comments (
  id text UNIQUE PRIMARY KEY,
  id_parents text NOT NULL,
  user_id text NOT NULL REFERENCES users(id),
  username text NOT NULL REFERENCES users(username),
  body text NOT NULL,
  body_html text NOT NULL,
  removed boolean NOT NULL DEFAULT false,
  updated_at timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
  created_at timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
  upvotes integer NOT NULL DEFAULT 0,
  downvotes integer NOT NULL DEFAULT 0
)
