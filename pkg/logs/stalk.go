package logs

import (
	"log"

	"github.com/kil0meters/acolyte/pkg/authorization"
	"github.com/kil0meters/acolyte/pkg/database"
)

// StalkUser fuzzy searches logs for a given query
func StalkUser(username string) []LogMessage {
	account := authorization.AccountFromUsername(username)
	if account == nil {
		return nil
	}

	rows, err := database.DB.Queryx("SELECT * FROM acolyte.chat_log WHERE account_id = $1 ORDER BY time DESC", account.ID)
	if err != nil {
		log.Println(err)
		return nil
	}

	resultRows := make([]LogMessage, 0)
	row := LogMessage{}

	for rows.Next() {
		rows.StructScan(&row)
		resultRows = append(resultRows, row)
	}

	return resultRows
}
