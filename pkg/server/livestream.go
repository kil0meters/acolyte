package server

import (
	"github.com/kil0meters/acolyte/pkg/authorization"
	"github.com/kil0meters/acolyte/pkg/homepage"
	"log"
	"net/http"
)

// ServeLivestream serves livestream page
func ServeLivestream(w http.ResponseWriter, _ *http.Request) {
	err := webTemplate.ExecuteTemplate(w, "livestream", struct {
		ChannelID string
	}{
		ChannelID: homepage.YoutubeChannelID,
	})
	if err != nil {
		log.Println(err)
	}
}

// ServeChat serves chat embed
func ServeChat(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	account := authorization.GetAccount(w, r)

	if !account.Permissions.AtLeast(authorization.LoggedOut) {
		http.Error(w, "hey banned buckaroo try not being banned before you chat :)", http.StatusUnauthorized)
		return
	}

	err := webTemplate.ExecuteTemplate(w, "chat", struct {
		Account       *authorization.Account
		IsModerator   bool
		IsStreamEmbed bool
	}{
		Account:       account,
		IsModerator:   account.Permissions.AtLeast(authorization.Moderator),
		IsStreamEmbed: r.Form.Get("stream_embed") == "1",
	})

	if err != nil {
		log.Println(err)
	}
}
