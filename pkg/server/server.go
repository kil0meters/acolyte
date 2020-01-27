package server

import (
	"net/http"

	"log"

	"github.com/kil0meters/acolyte/pkg/chat"
	"github.com/kil0meters/acolyte/pkg/homepage"

	"github.com/urfave/negroni"

	// "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	// "github.com/gorilla/websocket"
)

// StartServer starts the server
func StartServer() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1/").Subrouter()

	// live chat socket
	pool := chat.NewPool()
	go pool.Start()

	r.HandleFunc("/", homepage.ServeHomepage)

	r.PathPrefix("/scripts/").Handler(http.StripPrefix("/scripts/", http.FileServer(http.Dir("./acolyte-web/scripts/"))))
	r.PathPrefix("/styles/").Handler(http.StripPrefix("/styles/", http.FileServer(http.Dir("./acolyte-web/styles/"))))

	r.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "acolyte-web/chat.html")
	})

	api.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWS(pool, w, r)
	})

	log.Println("Starting server at http://localhost:3000")

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(r)

	http.ListenAndServe(":3000", n)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "acolyte-web/index.html")
	return
}
