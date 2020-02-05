package forum

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kil0meters/acolyte/pkg/authorization"
	"github.com/kil0meters/acolyte/pkg/database"
)

var frontpageTemplate *template.Template = template.Must(template.ParseFiles("./templates/forum/frontpage.html"))
var postEditorTemplate *template.Template = template.Must(template.ParseFiles("./templates/forum/post-editor.html"))

// Data contains data for the forum pages wow
type Data struct {
	IsLoggedIn bool
	Post       *Post
	Posts      []Post
	Account    *authorization.Account
}

// ServeForum serves forum front page
func ServeForum(w http.ResponseWriter, r *http.Request) {
	posts := []Post{}

	err := database.DB.Select(&posts, "SELECT * FROM acolyte.posts")
	if err != nil {
		log.Println(err)
	}

	isLoggedIn := false
	if user := authorization.IsAuthorized(w, r, authorization.Banned); user != nil {
		isLoggedIn = true
	}

	data := Data{
		IsLoggedIn: isLoggedIn,
		Posts:      posts,
	}

	frontpageTemplate.Execute(w, data)
}

// ServePostEditor serves the post editor
func ServePostEditor(w http.ResponseWriter, r *http.Request) {
	user := authorization.IsAuthorized(w, r, authorization.Banned)
	if user == nil {
		http.Redirect(w, r, "/log-in?target=/forum/create-post", 200)
	}
	postEditorTemplate.Execute(w, nil)
}

// CreatePostForm creates a new post
func CreatePostForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	title := html.EscapeString(strings.Trim(r.Form.Get("title"), " \n\t"))
	body := html.EscapeString(strings.Trim(r.Form.Get("body"), " \n\t"))
	link := html.EscapeString(strings.Trim(r.Form.Get("link"), " \n\t"))

	account := authorization.IsAuthorized(w, r, authorization.Standard)
	if account == nil {
		return
	}

	post, err := CreateNewPost(title, account, body, link)
	if err != nil {
		log.Println(err) // TODO: Unhandled error
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/forum/posts/%s", post.ID), http.StatusSeeOther)
}

// ServePost serves a post
func ServePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	log.Println(params["message_id"])

	post := PostFromID(params["message_id"])
	if post != nil {
		log.Println(post.Body)

		account := authorization.AccountFromID(post.AccountID)

		data := Data{
			IsLoggedIn: false,
			Post:       post,
			Account:    account,
		}

		postTemplate.Execute(w, data)
	}
}
