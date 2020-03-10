CREATE TABLE IF NOT EXISTS blog_posts (
  id text UNIQUE PRIMARY KEY,
  title text NOT NULL,
  body text NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
  created_at timestamp NOT NULL DEFAULT (now() at time zone 'utc')
)
