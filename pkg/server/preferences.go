package server

import (
	"github.com/kil0meters/acolyte/pkg/authorization"
	"log"
	"net/http"
)

func ServePreferencesHomepage(w http.ResponseWriter, r *http.Request) {
	account := authorization.GetAccount(w, r)

	if !account.Permissions.AtLeast(authorization.Standard) {
		http.Error(w, "hey banned buckaroo try not being banned before you chat :)", http.StatusUnauthorized)
		return
	}

	err := webTemplate.ExecuteTemplate(w, "preferences-home", nil)

	if err != nil {
		log.Println(err)
	}
}
