package authorization

import (
	"log"
	"net/http"

	"github.com/kil0meters/acolyte/pkg/database"
)

// PermissionLevel enum for different forum permissions
type PermissionLevel string

const (
	// Admin legends
	Admin PermissionLevel = "AUTH_ADMIN"
	// Moderator chat moderators
	Moderator PermissionLevel = "AUTH_MODERATOR"
	// Standard plebs
	Standard PermissionLevel = "AUTH_STANDARD"
	// Banned banned plebs
	Banned PermissionLevel = "AUTH_BANNED"
)

// IsAuthorized tests if a session token is valid
func IsAuthorized(r *http.Request, requiredPermission PermissionLevel) *Account {
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

	return nil
}
