package homepage

import (
	"net/http"
	"text/template"
)

var homepageTemplate *template.Template = template.Must(template.ParseFiles("./templates/homepage.html"))

// Data data for home page
type Data struct {
	FeaturedVideo YoutubeVideo
	IsLive        bool
	Header        []HeaderListElement
}

// HeaderListElement a single list item in the header
type HeaderListElement struct {
	Name string
	URL  string
}

// ServeHomepage serves the homepage
func ServeHomepage(w http.ResponseWriter, r *http.Request) {
	data := Data{
		FeaturedVideo: YoutubeVideo{
			Title:     "This is a test video",
			ID:        "g15-lvmIrcg",
			Thumbnail: "https://i.ytimg.com/vi/g15-lvmIrcg/hq720.jpg",
		},
		Header: []HeaderListElement{
			{Name: "Forum", URL: ""},
			{Name: "Videos", URL: ""},
			{Name: "Live", URL: "/chat"},
			{Name: "Blog", URL: ""},
			{Name: "Resume", URL: ""},
		},
		IsLive: true,
	}
	homepageTemplate.Execute(w, data)
}
