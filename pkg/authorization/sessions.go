package authorization

import (
	"encoding/gob"
	"net/http"
	"os"
	"time"

	"github.com/kil0meters/acolyte/pkg/database"

	"github.com/antonlindstrom/pgstore"
)

var store *pgstore.PGStore

// InitializeSessionManager initializes a session manager
func InitializeSessionManager() {
	store, _ = pgstore.NewPGStoreFromPool(database.DB.DB, []byte(os.Getenv("SECRET_KEY")))

	gob.Register(&Account{})

	store.Cleanup(time.Minute * 5)
}

// GetSession gets an http session
func GetSession(w http.ResponseWriter, r *http.Request) *Account {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	account := &Account{}
	accountInterface := session.Values["account"]
	if accountInterface == nil {
		account.Permissions = LoggedOut
	} else {
		account = accountInterface.(*Account)
	}

	return account
}

// CreateSession creates a new session cookie
func CreateSession(w http.ResponseWriter, r *http.Request, account *Account) bool {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	session.Values["account"] = account

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	return true
}

// InvalidateSession invalidates a session cookie
func InvalidateSession(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["account"] = Account{}
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
