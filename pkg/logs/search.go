package logs

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kil0meters/acolyte/pkg/database"
)

// SearchLogs fuzzy searches logs for a given query
func SearchLogs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	log.Println(r.URL)
	search := params["search"]
	log.Println(search)
	// username := params["from"]

	rows, err := database.DB.Queryx("SELECT * FROM acolyte.chat_log WHERE SIMILARITY(message, $1) > 0 ORDER BY SIMILARITY(message, $1) DESC LIMIT 100", search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resultRows := make([]LogMessage, 0)
	row := LogMessage{}

	for rows.Next() {
		rows.StructScan(&row)
		resultRows = append(resultRows, row)
	}

	log.Println(resultRows)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultRows)
}
