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
	var err error
	DB, err = sqlx.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}

	if err = DB.Ping(); err != nil {
		log.Panic(err)
	}

	DB.MustExec("CREATE SCHEMA IF NOT EXISTS acolyte")

	DB.MustExec(`CREATE TABLE IF NOT EXISTS acolyte.posts (id text UNIQUE PRIMARY KEY,
														   user_id text,
														   title text NOT NULL,
														   body text,
														   link text,
														   removed boolean DEFAULT FALSE,
														   created_at timestamp DEFAULT NOW(),
														   upvotes integer DEFAULT 0,
														   downvotes integer DEFAULT 0)`)

	DB.MustExec(`CREATE TABLE IF NOT EXISTS acolyte.accounts (id text UNIQUE PRIMARY KEY,
															  username text UNIQUE NOT NULL,
															  email text UNIQUE,
															  password_hash text NOT NULL,
															  created_at timestamp DEFAULT NOW(),
															  sessions text[] DEFAULT '{}'::text[])`)

	DB.MustExec(`CREATE TABLE IF NOT EXISTS acolyte.bans (user_id text UNIQUE PRIMARY KEY,
														  banned_until timestamp NOT NULL,
														  ban_reason text,
														  banned_by text NOT NULL)`)
}
