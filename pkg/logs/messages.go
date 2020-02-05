package logs

import (
	"log"
	"time"

	"github.com/kil0meters/acolyte/pkg/database"
)

// GetByDay returns all messages on a given day
func GetByDay(timestamp time.Time) []LogMessage {
	rows, err := database.DB.Queryx("SELECT * FROM acolyte.chat_log WHERE date_trunc('day', time) = $1 ORDER BY time DESC", timestamp)
	if err != nil {
		log.Panicln(err)
	}

	resultRows := make([]LogMessage, 0)
	row := LogMessage{}

	for rows.Next() {
		rows.StructScan(&row)
		resultRows = append(resultRows, row)
	}

	return resultRows
}
