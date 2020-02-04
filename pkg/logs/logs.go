package logs

import (
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var homepageTemplate *template.Template = template.Must(template.ParseFiles("./templates/logs/homepage.html"))
var searchTemplate *template.Template = template.Must(template.ParseFiles("./templates/logs/search.html"))
var stalkTemplate *template.Template = template.Must(template.ParseFiles("./templates/logs/stalk.html"))
var messagesTemplate *template.Template = template.Must(template.ParseFiles("./templates/logs/messages.html"))

// LogMessage a struct representing a message log
type LogMessage struct {
	MessageID uuid.UUID `json:"message_id" db:"message_id"`
	AccountID string    `json:"account_id"  db:"account_id"`
	Username  string    `json:"username"   db:"username"`
	Timestamp time.Time `json:"time"       db:"time"`
	Message   string    `json:"message"     db:"message"`
}

// ServeHomepage serves logs homepage
func ServeHomepage(w http.ResponseWriter, r *http.Request) {
	homepageTemplate.Execute(w, nil)
}

// ServeSearch serves logs search page
func ServeSearch(w http.ResponseWriter, r *http.Request) {
	searchTemplate.Execute(w, nil)
}

// ServeStalk serves stalk page
func ServeStalk(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	messages := StalkUser(params["username"])

	stalkTemplate.Execute(w, messages)
}

// ServeMessages serves messages
func ServeMessages(w http.ResponseWriter, r *http.Request) {
	messagesTemplate.Execute(w, nil)
}
