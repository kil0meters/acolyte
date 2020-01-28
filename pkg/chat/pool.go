package chat

import (
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
	for client := range pool.Clients {
		client.Conn.WriteJSON(message.Data)
	}
}

// Start starts a pool
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true

			messageData := MessageData{
				Username: "Server",
				ID:       uuid.Must(uuid.NewRandom()),
				Text:     "A user joined",
			}

			pool.BroadcastMessage(Message{
				Type: 1,
				Data: messageData,
			})
		case client := <-pool.Unregister:
			delete(pool.Clients, client)

			messageData := MessageData{
				Username: "Server",
				ID:       uuid.Must(uuid.NewRandom()),
				Text:     "A user left",
			}

			pool.BroadcastMessage(Message{
				Type: 1,
				Data: messageData,
			})
		case message := <-pool.Broadcast:
			pool.BroadcastMessage(message)
		}
	}
}
