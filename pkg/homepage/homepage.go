package homepage

import (
	"log"
	"net/http"
	"text/template"
	"time"
)

var homepageTemplate *template.Template = template.Must(template.ParseFiles("./templates/homepage.html"))
var isLive bool = false

// YoutubeChannelID id of the channel featured
// my channel - UCdXFe8CHwhS2nUT8JB5K2kQ
// chill beats - UCSJ4gkVC6NrvII8umztf0Ow
const YoutubeChannelID = "UCSJ4gkVC6NrvII8umztf0Ow"

// Data data for home page
type Data struct {
	FeaturedVideo YoutubeVideo
	ChannelID     string
	IsLive        bool
	Header        []HeaderListElement
}

// HeaderListElement a single list item in the header
type HeaderListElement struct {
	Name string
	URL  string
}

func checkIfLive() {
	_isLive := CheckIfChannelIsLive(YoutubeChannelID)
	if !_isLive && isLive {
		log.Println("Channel is no longer live")
	}
	if _isLive && !isLive {
		log.Println("Channel is now live")
	}

	isLive = _isLive
}

// CheckIfLiveJob Checks if livestreaming every 5 minutes
func CheckIfLiveJob() {
	checkIfLive()
	ticker := time.NewTicker(5 * time.Minute)
	go func(ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				checkIfLive()
			}
		}
	}(ticker)
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
			{Name: "Forum", URL: "/forum"},
			{Name: "Videos", URL: "https://youtube.com/channel/" + YoutubeChannelID},
			{Name: "Live", URL: "/live"},
			{Name: "Logs", URL: "/logs"},
			{Name: "Blog", URL: "/blog"},
			{Name: "Resume", URL: "/resume.pdf"},
		},
		ChannelID: YoutubeChannelID,
		IsLive:    isLive,
	}
	homepageTemplate.Execute(w, data)
}
