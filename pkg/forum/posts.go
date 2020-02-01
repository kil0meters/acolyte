package forum

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/kil0meters/acolyte/pkg/database"

	"github.com/asaskevich/govalidator"
)

var postTemplate *template.Template = template.Must(template.ParseFiles("./templates/forum/post.html"))

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
	UserID    string `db:"user_id"   valid:"ascii,required"`
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
func CreateNewPost(title string, user *User, body string, link string) (*Post, error) {
	post := Post{
		ID:     GenerateID(6),
		UserID: user.ID,
		Title:  title,
		Link:   link,
		Body:   body,
	}

	log.Println(post)

	if !post.IsValid() {
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

	return &post
}

// NewPost creates a new post
func NewPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	title := r.Form.Get("title")
	body := r.Form.Get("body")
	link := r.Form.Get("link")

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

	data := Data{
		IsLoggedIn: false,
		Post:       post,
	}

	postTemplate.Execute(w, data)
}
