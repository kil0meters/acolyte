package forum

import (
	"errors"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/kil0meters/acolyte/pkg/database"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidUserData shows invalid user data
var ErrInvalidUserData = errors.New("Received invalid user data")

// User struct represents a user
type User struct {
	ID       string         `db:"id"            valid:"alphanum"`
	Username string         `db:"username"      valid:"alphanum"`
	Email    string         `db:"email"         valid:"email"`
	Hash     string         `db:"password_hash" valid:"ascii"`
	Sessions pq.StringArray `db:"sessions"      valid:"-"`
}

// IsValid tests if a user struct contains valid data
func (user User) IsValid() bool {
	result, err := govalidator.ValidateStruct(user)

	if err != nil || result == false {
		return false
	}

	return true
}

func hashString(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func verifyHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateUser creates a new user struct
func CreateUser(username string, email string, password string) (*User, error) {
	hash, err := hashString(password)
	if err != nil {
		return nil, err
	}

	user := User{
		ID:       GenerateID(6),
		Username: username,
		Email:    email,
		Hash:     hash,
	}

	log.Printf("Creating user %s with username %s and email %s", user.ID, user.Username, user.Email)

	if !user.IsValid() {
		return nil, ErrInvalidUserData
	}

	_, err = database.DB.NamedExec("INSERT INTO acolyte.accounts (id, username, email, password_hash) VALUES (:id, :username, :email, :password_hash)", user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateSessionCookies creates a new session cookie
func CreateSessionCookies(w http.ResponseWriter, username string, password string) bool {
	sessionID := GenerateID(16)
	hashedSessionID, _ := hashString(sessionID)

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

// ServeLogin shows login screen
func ServeLogin(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")

	if target == "" {
		target = "/?login_success=1"
	}

	loginTemplate.Execute(w, target)
}

// LoginForm wow
func LoginForm(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Referer())
	r.ParseForm()

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	target := r.Form.Get("target")
	if target == "" {
		target = "/forum"
	}

	if CreateSessionCookies(w, username, password) {
		http.Redirect(w, r, target, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "log-in?error=1", http.StatusSeeOther)
	}
}

// ServeSignup shows signin screen
func ServeSignup(w http.ResponseWriter, r *http.Request) {
	signupTemplate.Execute(w, nil)
}

// SignupForm wow
func SignupForm(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Referer())
	r.ParseForm()

	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	_, err := CreateUser(username, email, password)
	if err == ErrInvalidUserData {
		http.Error(w, "Invalid user data", 400)
		log.Println(err)
	} else if err != nil {
		http.Error(w, "User with that email address or username already exists", 400)
		log.Println(err)
	} else {
		// sets session cookie because storing username/password pair is not safe
		CreateSessionCookies(w, username, password)
	}
}
