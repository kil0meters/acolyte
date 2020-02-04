package chat

import (
	"fmt"

	"github.com/google/uuid"
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
		for client := range pool.Clients {
			client.Conn.WriteJSON(message.Data)
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
				if client.Username == testClient.Username && client != testClient {
					hasAnotherSession = true
				}
			}

			if hasAnotherSession == false && client.Username != "ANON" {
				messageData := MessageData{
					Username: "Server",
					ID:       uuid.Must(uuid.NewRandom()),
					Text:     fmt.Sprintf("%s joined", client.Username),
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
				if client.Username == testClient.Username {
					hasAnotherSession = true
				}
			}

			if hasAnotherSession == false && client.Username != "ANON" {
				messageData := MessageData{
					Username: "Server",
					ID:       uuid.Must(uuid.NewRandom()),
					Text:     fmt.Sprintf("%s left", client.Username),
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
