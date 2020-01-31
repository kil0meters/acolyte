package chat

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client struct for chat client
type Client struct {
	Username string
	Conn     *websocket.Conn
	Pool     *Pool
}

// Message struct for handling websocket messages
type Message struct {
	Type int         `json:"type"`
	Body string      `json:"body"`
	Data MessageData `json:"data"`
}

// MessageData a struct containing data for a message
type MessageData struct {
	Username string    `json:"username"`
	ID       uuid.UUID `json:"id"`
	Text     string    `json:"text"`
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

		messageData := MessageData{}
		json.Unmarshal(p, &messageData)
		messageData.Username = c.Username // force username to be forum username
		messageData.Text = strings.TrimSpace(messageData.Text)
		messageData.ID = uuid.Must(uuid.NewRandom())

		message := Message{
			Type: messageType,
			Data: messageData,
		}

		c.Pool.Broadcast <- message
		log.Printf("[chat] [%s] <%s> %s\n", c.Conn.RemoteAddr(), message.Data.Username, message.Data.Text)
	}
}
