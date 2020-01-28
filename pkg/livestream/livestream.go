package livestream

import (
	"net/http"
	"text/template"

	"github.com/kil0meters/acolyte/pkg/homepage"
)

var livestreamTemplate *template.Template = template.Must(template.ParseFiles("./templates/livestream.html"))

// Data data for livestream page
type Data struct {
	ChannelID string
}

// ServeLivestream serves livestream page
func ServeLivestream(w http.ResponseWriter, r *http.Request) {
	data := Data{
		ChannelID: homepage.YoutubeChannelID,
	}

	livestreamTemplate.Execute(w, data)
}
