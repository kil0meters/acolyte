package forum

import (
	"html/template"
	"log"
	"net/http"

	"github.com/kil0meters/acolyte/pkg/database"
)

var frontpageTemplate *template.Template = template.Must(template.ParseFiles("./templates/forum/frontpage.html"))

// Data contains data for the forum pages wow
type Data struct {
	IsLoggedIn bool
	Posts      []Post
}

// ServeForum serves forum front page
func ServeForum(w http.ResponseWriter, r *http.Request) {
	posts := []Post{}

	err := database.DB.Select(&posts, "SELECT * FROM acolyte.posts")
	if err != nil {
		log.Println(err)
	}

	data := Data{
		IsLoggedIn: true,
		Posts:      posts,
	}

	frontpageTemplate.Execute(w, data)
}
