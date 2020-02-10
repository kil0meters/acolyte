package chat

import (
	"github.com/kil0meters/acolyte/pkg/logs"
	"sync"
)

// Pool an echo pool of websocket connections
type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
	mu         sync.Mutex
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

// KillAllConnections kills all connections of a user with a specific username
func (pool *Pool) KillAllConnections(username string) {
	pool.mu.Lock()
	for client := range pool.Clients {
		if client.Account.Username == username {
			_ = client.Conn.Close()
			pool.Unregister <- client
		}
	}
	pool.mu.Unlock()
}

// BroadcastMessage broadcasts a message to a given pool
func (pool *Pool) BroadcastMessage(message Message) {
	pool.mu.Lock()
	if message.Data.Username != "ANON" {
		logs.RecordMessage(message.Data.ID, message.Data.AccountID, message.Data.Username, message.Data.Text.(string))

		for client := range pool.Clients {
			_ = client.Write(message.Data)
		}
	}
	pool.mu.Unlock()
}

// GetUserList returns an array of usernames
func (pool *Pool) GetUserList() []string {
	pool.mu.Lock()
	usernames := make([]string, 0)

	for client := range pool.Clients {
		usernames = append(usernames, client.Account.Username)
	}

	pool.mu.Unlock()
	return usernames
}

// Start starts a pool
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.mu.Lock()
			pool.Clients[client] = true

			hasAnotherSession := false
			for testClient := range pool.Clients {
				if client.Account.Username == testClient.Account.Username && client != testClient {
					hasAnotherSession = true
				}
			}

			if hasAnotherSession == false && client.Account.Username != "ANON" {

				messageData := MessageData{
					Username: "user-add",
					Text:     client.Account.Username,
				}

				pool.mu.Unlock()
				pool.BroadcastMessage(Message{
					Type: 1,
					Data: messageData,
				})
				pool.mu.Lock()

			}
			pool.mu.Unlock()

		case client := <-pool.Unregister:
			pool.mu.Lock()
			delete(pool.Clients, client)

			hasAnotherSession := false
			for testClient := range pool.Clients {
				if client.Account.Username == testClient.Account.Username {
					hasAnotherSession = true
				}
			}

			if hasAnotherSession == false && client.Account.Username != "ANON" {
				messageData := MessageData{
					Username: "user-remove",
					Text:     client.Account.Username,
				}

				pool.mu.Unlock()
				pool.BroadcastMessage(Message{
					Type: 1,
					Data: messageData,
				})
				pool.mu.Lock()
			}
			pool.mu.Unlock()
		case message := <-pool.Broadcast:
			pool.BroadcastMessage(message)
		}
	}
}
