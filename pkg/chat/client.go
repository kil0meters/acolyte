package chat

import (
	"encoding/json"
	"html"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/kil0meters/acolyte/pkg/authorization"
)

// Client struct for chat client
type Client struct {
	Account *authorization.Account
	Conn    *websocket.Conn
	Pool    *Pool
}

// Message struct for handling websocket messages
type Message struct {
	Type int         `json:"type"`
	Body string      `json:"body"`
	Data MessageData `json:"data"`
}

// MessageData a struct containing data for a message
type MessageData struct {
	Username  string `json:"username"`
	AccountID string
	ID        uuid.UUID `json:"id"`
	Text      string    `json:"text"`
}

// ReadMessage reads a message from websocket data
func ReadMessage(messageType int, body []byte) Message {
	messageData := MessageData{}
	json.Unmarshal(body, &messageData)
	messageData.Text = strings.TrimSpace(messageData.Text)
	messageData.ID = uuid.Must(uuid.NewRandom())

	message := Message{
		Type: messageType,
		Data: messageData,
	}

	message.Data.Text = html.EscapeString(message.Data.Text)

	return message
}

// Read Reads messages from a given client
func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	if c.Account.Username == "ANON" {
		c.Conn.WriteMessage(websocket.TextMessage, []byte("UNAUTHORIZED"))
	}

	for {
		messageType, body, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}

		message := ReadMessage(messageType, body)
		message.Data.Username = c.Account.Username // force username to be forum username
		message.Data.AccountID = c.Account.ID

		c.Pool.Broadcast <- message
		log.Printf("[chat] [%s] <%s> %s\n", c.Conn.RemoteAddr(), message.Data.Username, message.Data.Text)
	}
}
