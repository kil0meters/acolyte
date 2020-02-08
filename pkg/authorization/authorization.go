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
	// LoggedOut supa plebs
	LoggedOut PermissionLevel = "AUTH_LOGGED_OUT"
	// Banned banned plebs
	Banned PermissionLevel = "AUTH_BANNED"
)

// AtLeast tests if a PermissionLevel is at least another PermissionLevel
func (permissions PermissionLevel) AtLeast(minimumPermission PermissionLevel) bool {
	if minimumPermission == Admin {
		return permissions == Admin
	} else if minimumPermission == Moderator {
		return permissions == Moderator || permissions == Admin
	} else if minimumPermission == Standard {
		return permissions == Standard || permissions == Moderator || permissions == Admin
	} else if minimumPermission == LoggedOut {
		return permissions == LoggedOut || permissions == Standard || permissions == Moderator || permissions == Admin
	}
	return true
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
	r.ParseForm()

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	account := AccountFromLogin(username, password)

	target := r.Form.Get("target")
	if target == "" {
		target = "/forum?logged_in=1"
	}

	if account != nil {
		CreateSession(w, r, account)
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

	account, err := CreateAccount(username, email, password)
	if err != nil {
		http.Redirect(w, r, "/sign-up?error=1", http.StatusSeeOther)
		log.Println(err)
		return
	}

	CreateSession(w, r, account)
	http.Redirect(w, r, target, http.StatusSeeOther)

	log.Println(account)
}
