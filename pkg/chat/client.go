package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client struct for chat client
type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

// Message struct for handling websocket messages
type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

// Read Reads messages from a given client
func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}
		message := Message{Type: messageType, Body: string(p)}

		c.Pool.Broadcast <- message
		log.Println("Received message", message)
	}
}
