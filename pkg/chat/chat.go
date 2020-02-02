package chat

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"

	"github.com/kil0meters/acolyte/pkg/forum"
)

var chatTemplate *template.Template = template.Must(template.ParseFiles("./templates/chat.html"))

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // TODO: this is not safe
}

// ServeWS allows a user to join the live chat room
func ServeWS(pool *Pool, w http.ResponseWriter, r *http.Request) {
	log.Println("Starting WS session from user", r.RemoteAddr)

	user := forum.IsAuthorized(r)
	if user == nil {
		user = &forum.User{
			Username: "ANON",
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	client := &Client{
		Username: user.Username,
		Conn:     conn,
		Pool:     pool,
	}

	pool.Register <- client
	client.Read()
}

// ServeChat serves chat embed
func ServeChat(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	isStreamEmbed := r.Form.Get("stream_embed") == "1"

	log.Println(isStreamEmbed)

	chatTemplate.Execute(w, isStreamEmbed)
}
