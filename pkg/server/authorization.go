package server

import (
	"github.com/kil0meters/acolyte/pkg/authorization"
	"log"
	"net/http"
)

// ServeLogin shows login screen
func ServeLogin(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	err := r.URL.Query().Get("error") == "1"

	if target == "" {
		target = "/?login_success=1"
	}

	data := struct {
		Target string
		Error  bool
	}{
		Target: target,
		Error:  err,
	}

	_ = webTemplate.ExecuteTemplate(w, "login", data)
}

// LoginForm wow
func LoginForm(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	account := authorization.AccountFromLogin(username, password)

	target := r.Form.Get("target")
	if target == "" {
		target = "/forum?logged_in=1"
	}

	if account != nil {
		authorization.CreateSession(w, r, account)
		http.Redirect(w, r, target, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/log-in?error=1", http.StatusSeeOther)
	}
}

// ServeSignup shows signin screen
func ServeSignup(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	err := r.URL.Query().Get("error") == "1"

	if target == "" {
		target = "/?login_success=1"
	}

	data := struct {
		Target string
		Error  bool
	}{
		Target: target,
		Error:  err,
	}

	_ = webTemplate.ExecuteTemplate(w, "signup", data)
}

// SignupForm wow
func SignupForm(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	target := r.Form.Get("target")
	if target == "" {
		target = "/forum?account_created=1"
	}

	account, err := authorization.CreateAccount(username, email, password)
	if err != nil {
		http.Redirect(w, r, "/sign-up?error=1", http.StatusSeeOther)
		log.Println(err)
		return
	}

	authorization.CreateSession(w, r, account)
	http.Redirect(w, r, target, http.StatusSeeOther)
}
