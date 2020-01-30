package forum

import (
	"errors"
	"log"
	"net/http"

	"github.com/kil0meters/acolyte/pkg/database"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

// ErrInvalidPostData shows invalid post data
var ErrInvalidPostData = errors.New("Received invalid user data")

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
	ID        string `db:"id"        valid:"ascii,required"`
	Title     string `db:"title"     valid:"ascii,required"`
	Link      string `db:"link"      valid:"url"`
	Body      string `db:"body"      valid:"ascii"`
	Upvotes   int    `db:"upvotes"   valid:"-"`
	Downvotes int    `db:"downvotes" valid:"-"`
}

// IsValid tests if a post contains valid data
func (post Post) IsValid() bool {
	result, err := govalidator.ValidateStruct(post)

	if err != nil || result == false {
		return false
	}

	return true
}

// CreateNewPost adds a new post to the database
func CreateNewPost(title string, body string, link string) (*Post, error) {
	post := Post{
		ID:    GenerateID(6),
		Title: title,
		Link:  link,
		Body:  body,
	}

	if !post.IsValid() {
		return nil, ErrInvalidPostData
	}

	_, err := database.DB.NamedExec("INSERT INTO acolyte.posts (ID, title, body, link) VALUES (:id, :title, :body, :link)", post)
	if err != nil {
		return &post, err
	}

	return &post, nil
}

// FetchPosts fetches n posts from
func FetchPosts(sortingType SortingType, n int, offsetIndex int) {

}

// ListPosts lists posts
func ListPosts(w http.ResponseWriter, r *http.Request) {

}

// NewPost creates a new post
func NewPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	post, err := CreateNewPost(params["title"], params["body"], params["link"])
	if err != nil {
		log.Println(err)
	}

	log.Println(post)
}
