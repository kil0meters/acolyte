package livestream

import (
	"net/http"
	"text/template"
)

var livestreamTemplate *template.Template = template.Must(template.ParseFiles("./templates/livestream.html"))

// ServeLivestream serves livestream page
func ServeLivestream(w http.ResponseWriter, r *http.Request) {
	livestreamTemplate.Execute(w, nil)
}
