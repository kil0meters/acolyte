package forum

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/kil0meters/acolyte/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidUserData shows invalid user data
var ErrInvalidUserData = errors.New("Received invalid user data")

// User struct represents a user
type User struct {
	ID       string `db:"id"            valid:"alphanum"`
	Username string `db:"username"      valid:"alphanum"`
	Email    string `db:"email"         valid:"email"`
	Hash     string `db:"password_hash" valid:"ascii"`
}

// IsValid tests if a user struct contains valid data
func (user User) IsValid() bool {
	result, err := govalidator.ValidateStruct(user)

	if err != nil || result == false {
		return false
	}

	return true
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func verifyHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateUser creates a new user struct
func CreateUser(username string, email string, password string) (*User, error) {
	hash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := User{
		ID:       GenerateID(6),
		Username: username,
		Email:    email,
		Hash:     hash,
	}

	if !user.IsValid() {
		return nil, ErrInvalidUserData
	}

	_, err = database.DB.NamedExec("INSERT INTO TABLE acolyte.users (id, username, password) VALUES (:id, :username, :email, :password_hash)", user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
