package chat

import (
	"encoding/json"
	"html"
	"log"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"

	"github.com/kil0meters/acolyte/pkg/authorization"
)

// Client struct for chat client
type Client struct {
	Account *authorization.Account
	Session *sessions.Session
	Conn    *websocket.Conn
	Pool    *Pool
	mu      sync.Mutex
}

// Message struct for handling websocket messages
type Message struct {
	Type int         `json:"type"`
	Body string      `json:"body"`
	Data MessageData `json:"data"`
}

// MessageData a struct containing data for a message
type MessageData struct {
	Username  string    `json:"username"`
	AccountID string    `json:"-"`
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

// Write writes data to a client with a mutex
func (c *Client) Write(data interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Conn.WriteJSON(data)
}

// UpdateCommands resends command list
func (c *Client) UpdateCommands() {
	authorizedCommands := make([]*Command, 0)

	for command := range Commands {
		if c.Account.Permissions.AtLeast(command.RequiredPermission) {
			authorizedCommands = append(authorizedCommands, command)
		}
	}
	c.Write(authorizedCommands)
}

// Read Reads messages from a given client
func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	if c.Account.Username == "ANON" {
		c.Write(MessageData{
			Username: "Server",
			Text:     "UNAUTHORIZED",
		})
	}

	c.UpdateCommands()

	for {
		messageType, body, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}

		message := ReadMessage(messageType, body)
		message.Data.Username = c.Account.Username // force username to be forum username
		message.Data.AccountID = c.Account.ID

		if message.Data.Text == "" { // ignore zero length messages -- should be stopped in the client as well
			continue
		}

		if message.Data.Text[0] == '/' {
			output := ParseCommand(c, message.Data.Text)

			c.Write(MessageData{
				Username: "->",
				Text:     output,
			})

			log.Printf("[command] [%s] <%s> %s", c.Conn.RemoteAddr(), message.Data.Username, message.Data.Text)
		} else {
			c.Pool.Broadcast <- message
			log.Printf("[chat] [%s] <%s> %s\n", c.Conn.RemoteAddr(), message.Data.Username, message.Data.Text)
		}
	}
}
