package authorization

import (
	"log"
	"net/http"
	"text/template"
)

var loginTemplate *template.Template = template.Must(template.ParseFiles("./templates/forum/login.html"))
var signupTemplate *template.Template = template.Must(template.ParseFiles("./templates/forum/signup.html"))

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
	r.ParseForm()

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	target := r.Form.Get("target")
	if target == "" {
		target = "/forum?logged_in=1"
	}

	if CreateSessionCookies(w, username, password) {
		http.Redirect(w, r, target, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/log-in?error=1", http.StatusSeeOther)
	}
}

// ServeSignup shows signin screen
func ServeSignup(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")

	if target == "" {
		target = "/?login_success=1"
	}

	signupTemplate.Execute(w, target)
}

// SignupForm wow
func SignupForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	target := r.Form.Get("target")
	if target == "" {
		target = "/forum?account_created=1"
	}

	_, err := CreateAccount(username, email, password)
	if err == ErrInvalidAccountData {
		http.Error(w, "Invalid account data", 400)
		log.Println(err)
	} else if err != nil {
		http.Error(w, "User with that email address or username already exists", 400)
		log.Println(err)
	} else {
		// sets session cookie because storing username/password pair is not safe
		if CreateSessionCookies(w, username, password) {
			http.Redirect(w, r, target, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/sign-in?error=1", http.StatusSeeOther)
		}
	}
}
