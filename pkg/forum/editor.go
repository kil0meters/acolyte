package forum

import (
	"net/http"
	"text/template"
)

var editorTemplate *template.Template = template.Must(template.ParseFiles("./templates/forum/post-editor.html"))

// ServePostEditor serves the post editor
func ServePostEditor(w http.ResponseWriter, r *http.Request) {
	user := IsAuthorized(r)
	if user == nil {
		http.Redirect(w, r, "/log-in?target=chat", 200)
	}
	editorTemplate.Execute(w, nil)
}
