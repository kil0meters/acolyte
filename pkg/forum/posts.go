package forum

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/kil0meters/acolyte/pkg/database"

	"github.com/microcosm-cc/bluemonday"

	"github.com/asaskevich/govalidator"
)

var postTemplate *template.Template = template.Must(template.ParseFiles("./templates/forum/post.html"))

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
	ID        string    `db:"id"         valid:"printableascii,required"`
	UserID    string    `db:"user_id"    valid:"printableascii,required"`
	Title     string    `db:"title"      valid:"type(string),required"`
	Link      string    `db:"link"       valid:"printableascii,optional"`
	Body      string    `db:"body"       valid:"type(string),optional"`
	Removed   bool      `db:"removed"    valid:"type(bool),optional"`
	CreatedAt time.Time `db:"created_at" valid:"type(time.Time),optional"`
	Upvotes   int       `db:"upvotes"    valid:"-"`
	Downvotes int       `db:"downvotes"  valid:"-"`
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
func CreateNewPost(title string, user *User, body string, link string) (*Post, error) {
	post := Post{
		ID:     GenerateID(6),
		UserID: user.ID,
		Title:  title,
		Link:   link,
		Body:   body,
	}

	if !post.IsValid() { // TODO: post.Link still needs to be validated
		return nil, ErrInvalidPostData
	}

	_, err := database.DB.NamedExec("INSERT INTO acolyte.posts (id, user_id, title, body, link) VALUES (:id, :user_id, :title, :body, :link)", post)
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

// PostFromID retrieves a post from an ID
func PostFromID(id string) *Post {
	post := Post{}
	row := database.DB.QueryRowx("SELECT * FROM acolyte.posts WHERE id = $1", id)

	err := row.StructScan(&post)
	if err != nil {
		log.Println(err)
		return nil
	}

	fmt.Printf("%+v\n", post)

	return &post
}

// NewPost creates a new post
func NewPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	title := strings.Trim(r.Form.Get("title"), " \n\t")
	body := strings.Trim(r.Form.Get("body"), " \n\t")
	link := strings.Trim(r.Form.Get("link"), " \n\t")

	user := IsAuthorized(r)
	if user == nil {
		http.Error(w, "error: Unauthorized", http.StatusUnauthorized)
		return
	}

	post, err := CreateNewPost(title, user, body, link)
	if err != nil {
		log.Println(err) // TODO: Unhandled error
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/forum/posts/%s", post.ID), http.StatusSeeOther)

	// log.Println(post)
}

// ServePost serves a post
func ServePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	post := PostFromID(params["id"])
	post.Body = string(bluemonday.UGCPolicy().SanitizeBytes([]byte(post.Body)))

	user := UserFromUserID(post.UserID)

	data := Data{
		IsLoggedIn: false,
		Post:       post,
		User:       user,
	}

	postTemplate.Execute(w, data)
}
