package forum

import (
	"html"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kil0meters/acolyte/pkg/authorization"
	"github.com/kil0meters/acolyte/pkg/database"
)

var frontpageTemplate = template.Must(template.ParseFiles("./templates/forum/frontpage.html"))
var postEditorTemplate = template.Must(template.ParseFiles("./templates/forum/post-editor.html"))

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

	err := database.DB.Select(&posts, "SELECT * FROM posts")
	if err != nil {
		log.Println(err)
	}

	account := authorization.GetAccount(w, r)

	if !account.Permissions.AtLeast(authorization.LoggedOut) {
		http.Error(w, "Hey buddy banned users aren't allowed here :)", http.StatusUnauthorized)
		return
	}

	isLoggedIn := false
	if account.Permissions.AtLeast(authorization.Standard) {
		isLoggedIn = true
	}

	data := Data{
		IsLoggedIn: isLoggedIn,
		Posts:      posts,
	}

	_ = frontpageTemplate.Execute(w, data)
}

// ServePostEditor serves the post editor
func ServePostEditor(w http.ResponseWriter, r *http.Request) {
	account := authorization.GetAccount(w, r)

	if !account.Permissions.AtLeast(authorization.LoggedOut) {
		http.Error(w, "Banned users aren't allowed to post dummy :)", http.StatusUnauthorized)
	}

	if account.Permissions.AtLeast(authorization.Standard) {
		_ = postEditorTemplate.Execute(w, nil)
	} else {
		http.Redirect(w, r, "/log-in?target=/forum/create-post", http.StatusTemporaryRedirect)
	}
}

// CreatePostForm creates a new post
func CreatePostForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	title := html.EscapeString(strings.Trim(r.Form.Get("title"), " \n\t"))
	body := html.EscapeString(strings.Trim(r.Form.Get("body"), " \n\t"))
	link := html.EscapeString(strings.Trim(r.Form.Get("link"), " \n\t"))

	account := authorization.GetAccount(w, r)
	if !account.Permissions.AtLeast(authorization.Standard) {
		http.Redirect(w, r, "/forum/create-post?error=1", http.StatusSeeOther)
	}

	post, err := CreateNewPost(title, account, body, link)
	if err != nil {
		log.Println(err) // TODO: Unhandled error
		return
	}

	http.Redirect(w, r, "/forum/posts/"+post.ID, http.StatusSeeOther)
}

// ServePost serves a post
func ServePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	log.Println(params["message_id"])

	account := authorization.GetAccount(w, r)

	if !account.Permissions.AtLeast(authorization.LoggedOut) {
		http.Error(w, "why don't you not be banned before you try to look at memes :thinking:", http.StatusUnauthorized)
		return
	}

	post := PostFromID(params["message_id"])
	if post != nil {
		log.Println(post.Body)

		posterAccount := authorization.AccountFromID(post.AccountID)

		data := Data{
			IsLoggedIn: account.Permissions.AtLeast(authorization.Standard),
			Post:       post,
			Account:    posterAccount,
		}

		_ = postTemplate.Execute(w, data)
	}
}
