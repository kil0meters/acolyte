package authorization

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/kil0meters/acolyte/pkg/database"

	"github.com/asaskevich/govalidator"
	"github.com/lib/pq"
)

// ErrInvalidAccountData shows invalid account data
var ErrInvalidAccountData = errors.New("Received invalid account data")

// Account struct represents a account
type Account struct {
	ID          string          `db:"account_id"    valid:"alphanum"`
	Username    string          `db:"username"      valid:"alphanum"`
	Email       string          `db:"email"         valid:"email"`
	Hash        string          `db:"password_hash" valid:"printableascii"`
	CreatedAt   time.Time       `db:"created_at"    valid:"-"`
	Permissions PermissionLevel `db:"permissions"   valid:"-"`
	Sessions    pq.StringArray  `db:"sessions"      valid:"-"`
}

// IsValid tests if a account struct contains valid data
func (account Account) IsValid() bool {
	result, err := govalidator.ValidateStruct(account)

	if err != nil || result == false {
		return false
	}

	return true
}

// AccountFromID gets a account from an id
func AccountFromID(id string) *Account {
	account := Account{}

	row := database.DB.QueryRowx("SELECT * FROM acolyte.accounts WHERE account_id = $1", id)

	err := row.StructScan(&account)
	if err != nil {
		return nil
	}

	return &account
}

// AccountFromUsername gets a account from a username
func AccountFromUsername(username string) *Account {
	account := Account{}

	row := database.DB.QueryRowx("SELECT * FROM acolyte.accounts WHERE username = $1", username)

	err := row.StructScan(&account)
	if err != nil {
		return nil
	}

	return &account
}

// CreateAccount creates a new account struct
func CreateAccount(username string, email string, password string) (*Account, error) {
	hash, err := HashString(password)
	if err != nil {
		return nil, err
	}

	account := Account{
		ID:       GenerateID(6),
		Username: username,
		Email:    email,
		Hash:     hash,
	}

	log.Printf("Creating account %s with username %s and email %s", account.ID, account.Username, account.Email)

	if !account.IsValid() {
		return nil, ErrInvalidAccountData
	}

	_, err = database.DB.NamedExec("INSERT INTO acolyte.accounts (account_id, username, email, password_hash) VALUES (:account_id, :username, :email, :password_hash)", account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// CreateSessionCookies creates a new session cookie
func CreateSessionCookies(w http.ResponseWriter, username string, password string) bool {
	sessionID := GenerateID(16)
	hashedSessionID, _ := HashString(sessionID)

	meme, err := database.DB.Exec("UPDATE acolyte.accounts SET sessions = array_append(sessions, $1) WHERE username = $2", hashedSessionID, username)
	if err != nil {
		log.Println(err)
		return false
	}

	rowsAffected, _ := meme.RowsAffected()

	if rowsAffected == 0 {
		return false
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Value: sessionID,

		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: username,

		SameSite: http.SameSiteStrictMode,
	})

	return true
}

// InvalidateSessionCookie invalidates a session cookie
func InvalidateSessionCookie(username string, password string, sessionHash string) error {
	_, err := database.DB.Exec("UPDATE acolyte.accounts SET sessions = array_remove(sessions, ?)  WHERE username = ?", sessionHash, username)

	return err
}
