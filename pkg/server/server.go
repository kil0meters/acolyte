package server

import (
	"net/http"
	"os"

	"log"

	"github.com/kil0meters/acolyte/pkg/chat"
	"github.com/kil0meters/acolyte/pkg/database"
	"github.com/kil0meters/acolyte/pkg/forum"
	"github.com/kil0meters/acolyte/pkg/homepage"
	"github.com/kil0meters/acolyte/pkg/livestream"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// StartServer starts the server
func StartServer() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1/").Subrouter()

	// live chat socket
	pool := chat.NewPool()
	go pool.Start()

	homepage.CheckIfLiveJob()
	database.InitDatabase("postgres://kilometers@localhost:5432/kilometers?sslmode=disable")

	r.HandleFunc("/", homepage.ServeHomepage)

	r.PathPrefix("/scripts/").Handler(http.StripPrefix("/scripts/", http.FileServer(http.Dir("./acolyte-web/scripts/"))))
	r.PathPrefix("/styles/").Handler(http.StripPrefix("/styles/", http.FileServer(http.Dir("./acolyte-web/styles/"))))
	r.HandleFunc("/forum", forum.ServeForum)
	r.HandleFunc("/chat", chat.ServeChat)
	r.HandleFunc("/live", livestream.ServeLivestream)

	api.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWS(pool, w, r)
	})
	// api.HandleFunc("/getPost", _).Queries("id", "{id:[a-zA-Z]{6}}")
	api.HandleFunc("/list-posts", forum.ListPosts).Queries(
		"sorting-type", "{sorting-type:(?:hot|top|controversial|new)}",
		"amount", "{amount:(?:0?[1-9]|[12][0-9]|3[012])}",
		"start", "{start:[0-9]+}")

	api.HandleFunc("/new-post", forum.NewPost).Queries(
		"title", "{title}",
		"body", "{body}",
		"link", "{link}")
	// api.HandleFunc("/getComment", _).Queries("id", "{id:[a-zA-Z]{6}}")

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	log.Println("Starting server at http://localhost:" + port)

	http.ListenAndServe(":"+port, n)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "acolyte-web/index.html")
	return
}
