package server

import (
	"github.com/gorilla/mux"
	"github.com/kil0meters/acolyte/pkg/authorization"
	"github.com/kil0meters/acolyte/pkg/database"
	"github.com/kil0meters/acolyte/pkg/forum"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strings"
)

// ServeForum serves forum front page
func ServeForum(w http.ResponseWriter, r *http.Request) {
	var posts []forum.Post

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

	err = webTemplate.ExecuteTemplate(w, "forum-frontpage", struct {
		LoginStatus bool
		Account     *authorization.Account
		Posts       []forum.Post
	}{
		LoginStatus: isLoggedIn,
		Account:     account,
		Posts:       posts,
	})
	if err != nil {
		log.Println(err)
	}
}

// ServePostEditor serves the post editor
func ServePostEditor(w http.ResponseWriter, r *http.Request) {
	account := authorization.GetAccount(w, r)

	if !account.Permissions.AtLeast(authorization.LoggedOut) {
		http.Error(w, "Banned users aren't allowed to post dummy :)", http.StatusUnauthorized)
	}

	if account.Permissions.AtLeast(authorization.Standard) {
		_ = webTemplate.ExecuteTemplate(w, "post-editor", nil)
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

	post, err := forum.CreateNewPost(title, account, body, link)
	if err != nil {
		log.Println(err) // TODO: Unhandled error
		return
	}

	http.Redirect(w, r, "/forum/posts/"+post.ID, http.StatusSeeOther)
}

// ServePost serves a post
func ServePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	account := authorization.GetAccount(w, r)

	if !account.Permissions.AtLeast(authorization.LoggedOut) {
		http.Error(w, "why don't you not be banned before you try to look at memes :thinking:", http.StatusUnauthorized)
		return
	}

	post := forum.PostFromID(params["message_id"])
	if post == nil {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	posterAccount := authorization.AccountFromID(post.AccountID)

	err := webTemplate.ExecuteTemplate(w, "forum-post", struct {
		LoginStatus   bool
		Account       *authorization.Account
		Post          *forum.Post
		PosterAccount *authorization.Account
	}{
		Account:       account,
		LoginStatus:   account.Permissions.AtLeast(authorization.Standard),
		Post:          post,
		PosterAccount: posterAccount,
	})
	if err != nil {
		log.Println(err)
	}
}