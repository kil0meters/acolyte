package homepage

import (
	"sync"
	"text/template"
	"time"
)

const YoutubeChannelID = "UCSJ4gkVC6NrvII8umztf0Ow"

var homepageTemplate = template.Must(template.ParseFiles("./templates/home.gohtml"))

type ChannelData struct {
	FeaturedVideo YoutubeVideo
	ChannelID     string
	Header        []HeaderListElement
	LiveStatus    bool
	Mu            sync.Mutex
}

// HeaderListElement a single list item in the header
type HeaderListElement struct {
	Name string
	URL  string
}

// my channel - UCdXFe8CHwhS2nUT8JB5K2kQ
// chill beats - UCSJ4gkVC6NrvII8umztf0Ow
var Data = &ChannelData{
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
	ChannelID:  YoutubeChannelID,
	LiveStatus: false,
}

func checkIfLive() {
	_isLive := CheckIfChannelIsLive(YoutubeChannelID)

	Data.Mu.Lock()
	Data.LiveStatus = _isLive
	Data.Mu.Unlock()
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
