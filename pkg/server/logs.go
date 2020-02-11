package server

import (
	"github.com/gorilla/mux"
	"github.com/kil0meters/acolyte/pkg/logs"
	"net/http"
	"time"
)

// ServeLogsFrontpage serves logs homepage
func ServeLogsFrontpage(w http.ResponseWriter, _ *http.Request) {
	t1, _ := time.Parse("2006-01-02", "2020-01-01")
	t2 := time.Now()

	dates := make([]time.Time, 0)
	dates = append(dates, t1)

	for t1.Before(t2) {
		t1 = t1.AddDate(0, 0, 1)
		dates = append(dates, t1)
	}

	dates = dates[:len(dates)-2]
	_ = webTemplate.ExecuteTemplate(w, "logs-frontpage", dates)
}

// ServeSearch serves logs search page
func ServeLogsSearch(w http.ResponseWriter, _ *http.Request) {
	_ = webTemplate.ExecuteTemplate(w, "logs-search", nil)
}

// ServeStalk serves stalk page
func ServeStalk(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	messages := logs.StalkUser(params["username"])

	_ = webTemplate.ExecuteTemplate(w, "logs-stalk", messages)
}

// ServeMessagesByDate serves messages on a specific day
func ServeMessagesByDate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	date, _ := time.Parse("2006-01-02", params["date"])

	result := logs.LogResult{
		Title:   "Messages on " + date.Format("2006-01-02"),
		Date:    date,
		Results: logs.GetByDay(date),
	}

	_ = webTemplate.ExecuteTemplate(w, "logs-messages", result)
}
