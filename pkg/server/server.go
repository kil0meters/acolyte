package server

import (
	"github.com/kil0meters/acolyte/pkg/authorization"
	"github.com/kil0meters/acolyte/pkg/chat"
	"github.com/kil0meters/acolyte/pkg/homepage"
	"github.com/kil0meters/acolyte/pkg/logs"
	"net/http"
	"os"
	"text/template"

	"log"

	"github.com/kil0meters/acolyte/pkg/database"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var webTemplate = template.Must(template.ParseFiles("./templates/headers.gohtml",
	"./templates/forum.gohtml",
	"./templates/home.gohtml",
	"./templates/livestream.gohtml",
	"./templates/logs.gohtml",
	"./templates/auth.gohtml"))

// StartServer starts the server
func StartServer() {
	r := mux.NewRouter().StrictSlash(true)
	api := r.PathPrefix("/api/v1/").Subrouter()
	forumRouter := r.PathPrefix("/forum").Subrouter()
	logsRouter := r.PathPrefix("/logs").Subrouter()

	// live chat socket
	pool := chat.NewPool()
	go pool.Start()

	database.InitDatabase(os.Getenv("DATABASE_URL"))
	chat.InitializeCommands()
	authorization.InitializeSessionManager()

	homepage.CheckIfLiveJob()
	authorization.CheckBansJob()

	r.HandleFunc("/", ServeHomepage)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./acolyte-web/dist/"))))

	forumRouter.HandleFunc("", ServeForum)
	forumRouter.HandleFunc("/create-post", ServePostEditor).Methods("GET")
	forumRouter.HandleFunc("/create-post", CreatePostForm).Methods("POST")
	forumRouter.HandleFunc("/posts/{message_id:[a-zA-Z]{6}}", ServePost)

	r.HandleFunc("/chat", ServeChat)
	r.HandleFunc("/live", ServeLivestream)

	r.HandleFunc("/log-in", ServeLogin).Methods("GET")
	r.HandleFunc("/log-in", LoginForm).Methods("POST")
	r.HandleFunc("/sign-up", ServeSignup).Methods("GET")
	r.HandleFunc("/sign-up", SignupForm).Methods("POST")

	logsRouter.HandleFunc("", ServeLogsFrontpage)
	logsRouter.HandleFunc("/search", ServeLogsSearch)
	logsRouter.HandleFunc("/stalk", ServeStalk).Queries("username", "{username}")
	logsRouter.HandleFunc("/messages/{date:(?:[12]\\d{3}-(?:0[1-9]|1[0-2])-(?:0[1-9]|[12]\\d|3[01]))}", ServeMessagesByDate)

	api.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWS(pool, w, r)
	})

	api.HandleFunc("/search-logs", logs.SearchLogs).Queries(
		"search", "{search}")
	// "username", "{username}")

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	log.Println("Starting server at http://localhost:" + port)

	http.ListenAndServe(":"+port, n)
}
