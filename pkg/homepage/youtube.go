package homepage

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// YoutubeVideo a struct containing Data for a youtube video
type YoutubeVideo struct {
	Title     string
	ID        string
	Thumbnail string
}

// CheckIfChannelIsLive checks if a channel is live
func CheckIfChannelIsLive(channelID string) bool {
	res, err := http.Get("https://www.youtube.com/channel/" + channelID)
	if err != nil {
		log.Println("Encountered an error when attempting to check livestreaming status")
		return false
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Println("Encountered error", res.StatusCode)
	}

	document, err := goquery.NewDocumentFromReader(res.Body)

	_isLive := false

	document.Find(".yt-badge-live").First().Each(func(index int, selection *goquery.Selection) {
		_isLive = true
	})

	return _isLive
}
