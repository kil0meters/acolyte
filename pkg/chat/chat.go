package chat

import (
	"log"

	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // TODO: this is not safe
}

// ServeWS allows a user to join the live chat room
func ServeWS(pool *Pool, w http.ResponseWriter, r *http.Request) {
	log.Println("Starting WS session")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	client := &Client{
		Conn: conn,
		Pool: pool,
	}

	println(client)

	pool.Register <- client
	client.Read()
}
