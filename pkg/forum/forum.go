package forum

import (
	"html/template"
	"log"
	"net/http"

	"github.com/kil0meters/acolyte/pkg/authorization"
	"github.com/kil0meters/acolyte/pkg/database"
)

var frontpageTemplate *template.Template = template.Must(template.ParseFiles("./templates/forum/frontpage.html"))
var loginTemplate *template.Template = template.Must(template.ParseFiles("./templates/forum/login.html"))
var signupTemplate *template.Template = template.Must(template.ParseFiles("./templates/forum/signup.html"))

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
	if user := authorization.IsAuthorized(r, authorization.Banned); user != nil {
		isLoggedIn = true
	}

	data := Data{
		IsLoggedIn: isLoggedIn,
		Posts:      posts,
	}

	frontpageTemplate.Execute(w, data)
}
