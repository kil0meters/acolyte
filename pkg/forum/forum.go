package forum

import (
	"html/template"
	"log"
	"net/http"

	"github.com/kil0meters/acolyte/pkg/database"
)

var frontpageTemplate *template.Template = template.Must(template.ParseFiles("./templates/forum/frontpage.html"))
var loginTemplate *template.Template = template.Must(template.ParseFiles("./templates/forum/login.html"))
var signupTemplate *template.Template = template.Must(template.ParseFiles("./templates/forum/signup.html"))

// Data contains data for the forum pages wow
type Data struct {
	IsLoggedIn bool
	Posts      []Post
}

// IsAuthorized tests if a session token is valid
func IsAuthorized(r *http.Request) *User {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}

	usernameCookie, err := r.Cookie("username")
	if err != nil {
		return nil
	}

	user := User{}

	sessionToken := sessionCookie.Value
	username := usernameCookie.Value

	row := database.DB.QueryRowx("SELECT * FROM acolyte.accounts WHERE username = $1", username)
	if err != nil {
		log.Println(err)
		return nil
	}

	err = row.StructScan(&user)
	if err != nil {
		log.Println(err)
		return nil
	}

	isValid := false
	for i := 0; i < len(user.Sessions); i++ {
		if verifyHash(sessionToken, user.Sessions[i]) {
			isValid = true
		}
	}

	if isValid {
		return &user
	}

	return nil
}

// ServeForum serves forum front page
func ServeForum(w http.ResponseWriter, r *http.Request) {
	posts := []Post{}

	err := database.DB.Select(&posts, "SELECT * FROM acolyte.posts")
	if err != nil {
		log.Println(err)
	}

	isLoggedIn := false
	if user := IsAuthorized(r); user != nil {
		isLoggedIn = true
	}

	data := Data{
		IsLoggedIn: isLoggedIn,
		Posts:      posts,
	}

	frontpageTemplate.Execute(w, data)
}
