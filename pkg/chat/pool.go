package chat

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kil0meters/acolyte/pkg/logs"
)

// Pool an echo pool of websocket connections
type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

// NewPool creates a new pool object
func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

// BroadcastMessage broadcasts a message to a given pool
func (pool *Pool) BroadcastMessage(message Message) {
	if message.Data.Username != "ANON" {
		logs.RecordMessage(message.Data.ID, message.Data.AccountID, message.Data.Username, message.Data.Text)

		for client := range pool.Clients {
			client.Write(message.Data)
		}
	}
}

// Start starts a pool
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true

			hasAnotherSession := false
			for testClient := range pool.Clients {
				if client.Account.Username == testClient.Account.Username && client != testClient {
					hasAnotherSession = true
				}
			}

			if hasAnotherSession == false && client.Account.Username != "ANON" {
				messageData := MessageData{
					Username: "Server",
					ID:       uuid.Must(uuid.NewRandom()),
					Text:     fmt.Sprintf("%s joined", client.Account.Username),
				}

				pool.BroadcastMessage(Message{
					Type: 1,
					Data: messageData,
				})
			}
		case client := <-pool.Unregister:
			delete(pool.Clients, client)

			hasAnotherSession := false
			for testClient := range pool.Clients {
				if client.Account.Username == testClient.Account.Username {
					hasAnotherSession = true
				}
			}

			if hasAnotherSession == false && client.Account.Username != "ANON" {
				messageData := MessageData{
					Username: "Server",
					ID:       uuid.Must(uuid.NewRandom()),
					Text:     fmt.Sprintf("%s left", client.Account.Username),
				}

				pool.BroadcastMessage(Message{
					Type: 1,
					Data: messageData,
				})
			}
		case message := <-pool.Broadcast:
			pool.BroadcastMessage(message)
		}
	}
}
