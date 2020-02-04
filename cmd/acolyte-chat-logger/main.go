package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/kil0meters/acolyte/pkg/authorization"
	"github.com/kil0meters/acolyte/pkg/chat"
	"github.com/kil0meters/acolyte/pkg/database"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	database.InitDatabase(os.Getenv("DATABASE_URL"))

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("localhost:%s", port), Path: "/api/v1/chat"}
	log.Printf("Serving to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			messageType, body, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			if string(body) != "UNAUTHORIZED" {
				message := chat.ReadMessage(messageType, body)

				user := authorization.AccountFromUsername(message.Data.Username)

				if user != nil {
					log.Printf("%s [%s] %s", time.Now().Format("2006-01-02 15:04:05"), message.Data.Username, message.Data.Text)
					_, err = database.DB.Exec("INSERT INTO acolyte.chat_log (message_id, account_id, username, message) VALUES ($1, $2, $3, $4)", message.Data.ID, user.ID, user.Username, message.Data.Text)
					if err != nil {
						log.Panic(err)
					}
				}
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		// case t := <-ticker.C:
		// 	err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		// 	if err != nil {
		// 		log.Println("write:", err)
		// 		return
		// 	}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
