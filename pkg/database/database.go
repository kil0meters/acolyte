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

	_, err = DB.Exec("CREATE SCHEMA IF NOT EXISTS acolyte")
	if err != nil {
		log.Panic(err)
	}

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS acolyte.posts (id text UNIQUE PRIMARY KEY, user_id text, title text NOT NULL, body text, link text, upvotes integer DEFAULT 0, downvotes integer DEFAULT 0)")
	if err != nil {
		log.Panic(err)
	}

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS acolyte.accounts (id text UNIQUE PRIMARY KEY, username text UNIQUE NOT NULL, email text UNIQUE, password_hash text NOT NULL, sessions text[] DEFAULT '{}'::text[])")
}
