package authorization

import (
	"errors"
	"log"
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
	Hash        string          `db:"password_hash" valid:"printableascii"`
	CreatedAt   time.Time       `db:"created_at"    valid:"-"`
	Permissions PermissionLevel `db:"permissions"   valid:"-"`
	Sessions    pq.StringArray  `db:"sessions"      valid:"-"`
}

// IsValid tests if a account struct contains valid data
func (account Account) IsValid() bool {
	result, err := govalidator.ValidateStruct(account)

	if err != nil || result == false {
		log.Println(err)
		return false
	}

	return true
}

// Ban bans a user for time.Duration
func (account *Account) Ban(duration time.Duration) {
	account.Permissions = Banned

	database.DB.Exec("UPDATE accounts SET permissions = $1 WHERE account_id = $2", account.Permissions, account.ID)
}

// AccountFromID gets a account from an id
func AccountFromID(id string) *Account {
	account := Account{}

	row := database.DB.QueryRowx("SELECT * FROM accounts WHERE account_id = $1", id)

	err := row.StructScan(&account)
	if err != nil {
		return nil
	}

	return &account
}

// AccountFromUsername gets a account from a username
func AccountFromUsername(username string) *Account {
	account := Account{}

	row := database.DB.QueryRowx("SELECT * FROM accounts WHERE username = $1", username)

	err := row.StructScan(&account)
	if err != nil {
		return nil
	}

	return &account
}

// AccountFromLogin returns an account only if password matches hash
func AccountFromLogin(username string, password string) *Account {
	account := Account{}

	err := database.DB.QueryRowx("SELECT * FROM accounts WHERE username = $1", username).StructScan(&account)
	if err != nil {
		log.Println(err)
		return nil
	}

	if VerifyHash(password, account.Hash) {
		return &account
	}

	return nil
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
		Hash:     hash,
	}

	if !account.IsValid() {
		return nil, ErrInvalidAccountData
	}

	_, err = database.DB.NamedExec("INSERT INTO accounts (account_id, username, password_hash) VALUES (:account_id, :username, :password_hash)", account)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &account, nil
}
