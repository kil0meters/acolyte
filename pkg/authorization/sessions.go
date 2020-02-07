package authorization

import (
	"encoding/gob"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
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

// GetAccount gets an account object
func GetAccount(w http.ResponseWriter, r *http.Request) *Account {
	session := GetSession(w, r)
	if session == nil {
		return nil
	}

	account := &Account{}
	accountInterface := session.Values["account"]
	if accountInterface == nil {
		account.Permissions = LoggedOut
	} else {
		account = accountInterface.(*Account)

		updatedAccount := AccountFromID(account.ID)
		session.Values["account"] = updatedAccount

		err := store.Save(r, w, session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}
	}

	return account
}

// GetSession gets an http session
func GetSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return session
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
