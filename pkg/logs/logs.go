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
	AccountID string    `json:"account_id" db:"account_id"`
	Username  string    `json:"username"   db:"username"`
	Timestamp time.Time `json:"time"       db:"time"`
	Message   string    `json:"message"    db:"message"`
}

// LogResult contains result for a log request
type LogResult struct {
	Title   string
	Date    time.Time
	Results []LogMessage
}

// ServeHomepage serves logs homepage
func ServeHomepage(w http.ResponseWriter, r *http.Request) {
	t1, _ := time.Parse("2006-01-02", "2020-01-01")
	t2 := time.Now()

	dates := make([]time.Time, 0)
	dates = append(dates, t1)

	for t1.Before(t2) {
		t1 = t1.AddDate(0, 0, 1)
		dates = append(dates, t1)
	}

	dates = dates[:len(dates)-1]
	homepageTemplate.Execute(w, dates)
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

// ServeMessagesByDate serves messages on a specific day
func ServeMessagesByDate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	date, _ := time.Parse("2006-01-02", params["date"])

	result := LogResult{
		Title:   "Messages on " + date.Format("2006-01-02"),
		Date:    date,
		Results: GetByDay(date),
	}

	messagesTemplate.Execute(w, result)
}
