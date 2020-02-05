package authorization

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kil0meters/acolyte/pkg/database"
)

// IsAuthorized tests if a session token is valid
func IsAuthorized(w http.ResponseWriter, r *http.Request, requiredPermission PermissionLevel) *Account {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}

	usernameCookie, err := r.Cookie("username")
	if err != nil {
		return nil
	}

	account := Account{}

	sessionToken := sessionCookie.Value
	username := usernameCookie.Value

	row := database.DB.QueryRowx("SELECT * FROM acolyte.accounts WHERE username = $1", username)

	err = row.StructScan(&account)
	if err != nil {
		log.Println(err)
		return nil
	}

	isValid := false
	for i := 0; i < len(account.Sessions); i++ {
		if VerifyHash(sessionToken, account.Sessions[i]) {
			isValid = true
		}
	}

	if isValid {
		return &account
	}

	http.Error(w, fmt.Sprintf("You are not authorized to view this page. Required permission: %s but you are %s", requiredPermission, account.Permissions), http.StatusUnauthorized)

	return nil
}
