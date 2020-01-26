package chat

import "log"

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
func (pool *Pool) BroadcastMessage(message interface{}) {
	for client := range pool.Clients {
		log.Println(message)
		client.Conn.WriteJSON(message)
	}
}

// Start starts a pool
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			pool.BroadcastMessage(Message{Type: 1, Body: "A new user joined"})
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			pool.BroadcastMessage(Message{Type: 1, Body: "A user left"})
		case message := <-pool.Broadcast:
			pool.BroadcastMessage(message)
		}
	}
}
