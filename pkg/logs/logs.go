package logs

import (
	"html/template"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kil0meters/acolyte/pkg/database"
)

// var homepageTemplate = template.Must(template.ParseFiles("./templates/logs/home.gohtml"))
// var searchTemplate = template.Must(template.ParseFiles("./templates/logs/search.html"))
// var stalkTemplate = template.Must(template.ParseFiles("./templates/logs/stalk.html"))
// var messagesTemplate = template.Must(template.ParseFiles("./templates/logs/messages.html"))

// LogMessage a struct representing a message log
type LogMessage struct {
	MessageID uuid.UUID     `json:"message_id" db:"message_id"`
	AccountID string        `json:"account_id" db:"account_id"`
	Username  string        `json:"username"   db:"username"`
	Timestamp time.Time     `json:"time"       db:"time"`
	Message   template.HTML `json:"message"    db:"message"`
}

// LogResult contains result for a log request
type LogResult struct {
	Title   string
	Date    time.Time
	Results []LogMessage
}

// RecordMessage logs a message
func RecordMessage(messageID uuid.UUID, accountID string, username string, message string) {
	if accountID == "" {
		return
	}

	_, err := database.DB.Exec("INSERT INTO chat_log (message_id, account_id, username, message) VALUES ($1, $2, $3, $4)", messageID, accountID, username, message)
	if err != nil {
		log.Println("error logging message:", err.Error())
	}
}
