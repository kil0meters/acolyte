package server

import (
	"net/http"
	"os"

	"log"

	"github.com/kil0meters/acolyte/pkg/authorization"
	"github.com/kil0meters/acolyte/pkg/chat"
	"github.com/kil0meters/acolyte/pkg/database"
	"github.com/kil0meters/acolyte/pkg/forum"
	"github.com/kil0meters/acolyte/pkg/homepage"
	"github.com/kil0meters/acolyte/pkg/livestream"
	"github.com/kil0meters/acolyte/pkg/logs"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

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

	r.HandleFunc("/", homepage.ServeHomepage)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./acolyte-web/dist/"))))

	forumRouter.HandleFunc("", forum.ServeForum)
	forumRouter.HandleFunc("/create-post", forum.ServePostEditor).Methods("GET")
	forumRouter.HandleFunc("/create-post", forum.CreatePostForm).Methods("POST")
	forumRouter.HandleFunc("/posts/{message_id:[a-zA-Z]{6}}", forum.ServePost)

	r.HandleFunc("/chat", chat.ServeChat)
	r.HandleFunc("/live", livestream.ServeLivestream)

	r.HandleFunc("/log-in", authorization.ServeLogin).Methods("GET")
	r.HandleFunc("/log-in", authorization.LoginForm).Methods("POST")
	r.HandleFunc("/sign-up", authorization.ServeSignup).Methods("GET")
	r.HandleFunc("/sign-up", authorization.SignupForm).Methods("POST")

	logsRouter.HandleFunc("", logs.ServeHomepage)
	logsRouter.HandleFunc("/search", logs.ServeSearch)
	logsRouter.HandleFunc("/stalk", logs.ServeStalk).Queries("username", "{username}")
	logsRouter.HandleFunc("/messages/{date:(?:[12]\\d{3}-(?:0[1-9]|1[0-2])-(?:0[1-9]|[12]\\d|3[01]))}", logs.ServeMessagesByDate)
	// logsRouter.HandleFunc("/messages/{message_id}", logs.ServeLogs)

	api.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWS(pool, w, r)
	})
	api.HandleFunc("/list-posts", forum.ListPosts).Queries(
		"sorting-type", "{sorting-type:(?:hot|top|controversial|new)}",
		"amount", "{amount:(?:0?[1-9]|[12][0-9]|3[012])}",
		"start", "{start:[0-9]+}")

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
