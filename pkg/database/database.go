package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	// needs to be done I think
	_ "github.com/lib/pq"
)

// DB main database object
var DB *sqlx.DB

// InitDatabase initializes database
func InitDatabase(connStr string) {
	log.Println("Connecting to database:", connStr)

	var err error
	DB, err = sqlx.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}

	if err = DB.Ping(); err != nil {
		log.Panic(err)
	}

	// DB.MustExec("CREATE SCHEMA IF NOT EXISTS acolyte")
	// DB.MustExec("SET search_path TO acolyte,public")

	// this fails if the type doesn't exist initially..
	// DB.MustExec(`CREATE TYPE permission_level AS ENUM ('AUTH_ADMIN',
	// 																'AUTH_MODERATOR',
	// 																'AUTH_STANDARD',
	// 																'AUTH_BANNED')`)

	// We have to use (now() at time zone 'utc') because pq assumes all timezones are in that format

	DB.MustExec(`CREATE TABLE IF NOT EXISTS posts (post_id text UNIQUE PRIMARY KEY,
														account_id text,
														title text NOT NULL,
														body text,
														link text,
														removed boolean DEFAULT FALSE,
														created_at timestamp DEFAULT (now() at time zone 'utc'),
														upvotes integer DEFAULT 0,
														downvotes integer DEFAULT 0)`)

	// Uses extension citext
	DB.MustExec(`CREATE TABLE IF NOT EXISTS accounts (account_id text UNIQUE PRIMARY KEY,
															username citext UNIQUE NOT NULL,
															password_hash text NOT NULL,
															created_at timestamp DEFAULT (now() at time zone 'utc'),
															permissions permission_level DEFAULT 'AUTH_STANDARD')`)

	DB.MustExec(`CREATE TABLE IF NOT EXISTS bans (account_id text NOT NULL,
														unban_time timestamp NOT NULL,
														ban_time timestamp DEFAULT (now() at time zone 'utc'),
														ban_reason text,
														banned_by text NOT NULL)`)

	DB.MustExec(`CREATE TABLE IF NOT EXISTS chat_log (message_id uuid UNIQUE PRIMARY KEY,
															account_id text,
															username citext NOT NULL,
															time timestamp DEFAULT (now() at time zone 'utc'),
															message text NOT NULL)`)

	DB.MustExec(`CREATE TABLE IF NOT EXISTS comments (comment_id text UNIQUE PRIMARY KEY,
	                                                         parent_id text NOT NULL, 
	                                                         post_id text NOT NULL, 
	                                                         account_id text NOT NULL,
	                                                         username citext NOT NULL,
	                                                         created_at timestamp DEFAULT (now() at time zone 'utc'),
	                                                         body text NOT NULL,
	                                                         removed boolean DEFAULT false,
	                                                         upvotes integer DEFAULT 0,
	                                                         downvotes integer DEFAULT 0)`)

	DB.MustExec(`CREATE TABLE IF NOT EXISTS link_cache (link text UNIQUE PRIMARY KEY,
	                                                           title text,
	                                                           icon text,
	                                                           content text,
	                                                           created_at timestamp DEFAULT (now() at time zone 'utc'))`)
}
