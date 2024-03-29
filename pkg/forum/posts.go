package forum

import (
	"errors"
	"github.com/kil0meters/acolyte/pkg/links"
	"log"
	"time"

	"github.com/asaskevich/govalidator"

	"github.com/kil0meters/acolyte/pkg/authorization"
	"github.com/kil0meters/acolyte/pkg/database"
)

// ErrInvalidPostData shows invalid post data
var ErrInvalidPostData = errors.New("Received invalid post data")

// SortingType the type of sorting to use for posts/comments
type SortingType int

const (
	// Top sort posts by most upvotes
	Top SortingType = 1
	// New sort posts by time
	New SortingType = 2
	// Hot sort posts by magic fancy algorithm
	Hot SortingType = 3
	// Controversial sort posts by min(upvotes-downvotes)
	Controversial SortingType = 4
)

// Post struct containing data for a post
type Post struct {
	ID        string        `db:"post_id"    valid:"printableascii,required"`
	AccountID string        `db:"account_id" valid:"printableascii,required"`
	Username  string        `db:"username"   valid:"alphanum"`
	Title     string        `db:"title"      valid:"type(string),required"`
	LinkStr   string        `db:"link"       valid:"printableascii,optional"`
	Link      links.Article `db:"-"          valid:"-"`
	Body      string        `db:"body"       valid:"type(string),optional"`
	Removed   bool          `db:"removed"    valid:"-"`
	CreatedAt time.Time     `db:"created_at" valid:"-"`
	Upvotes   int           `db:"upvotes"    valid:"-"`
	Downvotes int           `db:"downvotes"  valid:"-"`
	Replies   []*Comment    `db:"-"          valid:"-"`
}

// IsValid tests if a post contains valid data
func (post Post) IsValid() bool {
	result, err := govalidator.ValidateStruct(post)

	if err != nil || result == false {
		log.Println(err)
		return false
	}

	return true
}

// CreateNewPost adds a new post to the database
func CreateNewPost(title string, account *authorization.Account, body string, link string) (*Post, error) {
	post := Post{
		ID:        authorization.GenerateID("p", 6),
		AccountID: account.ID,
		Username:  account.Username,
		Title:     title,
		LinkStr:   link,
		Body:      body,
	}

	if !post.IsValid() { // TODO: post.LinkStr still needs to be validated
		return nil, ErrInvalidPostData
	}

	_, err := database.DB.NamedExec("INSERT INTO posts (post_id, account_id, username, title, body, link) VALUES (:post_id, :account_id, :username, :title, :body, :link)", post)
	if err != nil {
		return &post, err
	}

	return &post, nil
}

// PostFromID retrieves a post from an ID
func PostFromID(id string) *Post {
	post := Post{}
	err := database.DB.QueryRowx("SELECT * FROM posts WHERE post_id = $1", id).StructScan(&post)

	if err != nil {
		log.Println(err)
		return nil
	}

	return &post
}

func PostsFromUsername(username string) []*Post {
	rows, err := database.DB.Queryx("SELECT * FROM posts WHERE username = $1 ORDER BY created_at DESC", username)
	if err != nil {
		log.Println(err)
		return nil
	}

	posts := make([]*Post, 0)

	for rows.Next() {
		post := new(Post)

		err = rows.StructScan(post)
		if err != nil {
			log.Println(err)
			return nil
		}

		posts = append(posts, post)
	}

	return posts
}
