package server

import (
	"github.com/kil0meters/acolyte/pkg/homepage"
	"log"
	"net/http"
)

// ServeHomepage serves the homepage
func ServeHomepage(w http.ResponseWriter, _ *http.Request) {
	err := webTemplate.ExecuteTemplate(w, "home", homepage.Data)
	if err != nil {
		log.Println(err)
	}
}
