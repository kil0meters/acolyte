package chat

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/kil0meters/acolyte/pkg/authorization"
)

var chatTemplate *template.Template = template.Must(template.ParseFiles("./templates/chat.html"))

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // TODO: this is not safe
}

// ServeWS allows a account to join the live chat room
func ServeWS(pool *Pool, w http.ResponseWriter, r *http.Request) {
	log.Println("Starting WS session from address", r.RemoteAddr)

	account := authorization.GetAccount(w, r)
	session := authorization.GetSession(w, r)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if !account.Permissions.AtLeast(authorization.LoggedOut) {
		ban, err := account.GetBanInfo()
		if err != nil {
			log.Println(err)
			conn.WriteJSON(MessageData{
				Username: "System",
				ID:       uuid.New(),
				Text:     fmt.Sprintf("Error receiving unban information :("),
			})
		} else {
			conn.WriteJSON(MessageData{
				Username: "System",
				ID:       uuid.New(),
				Text:     fmt.Sprintf("You are banned until %s", ban.UnbanTime.Format("January 2, 2006")),
			})
		}

		conn.Close()
	} else {
		client := &Client{
			Account: account,
			Session: session,
			Conn:    conn,
			Pool:    pool,
		}

		pool.Register <- client
		client.Read()
	}
}

// ServeChat serves chat embed
func ServeChat(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	isStreamEmbed := r.Form.Get("stream_embed") == "1"

	log.Println(isStreamEmbed)

	chatTemplate.Execute(w, isStreamEmbed)
}
