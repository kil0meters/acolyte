package server

import (
	"net/http"

	"log"

	"github.com/kil0meters/acolyte/pkg/chat"
	"github.com/urfave/negroni"
	// "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	// "github.com/gorilla/websocket"
)

// StartServer starts the server
func StartServer() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1/").Subrouter()

	r.HandleFunc("/", indexHandler)

	// live chat socket
	pool := chat.NewPool()
	go pool.Start()

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
