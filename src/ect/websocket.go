package ect

import (
	"time"
	"os/signal"
	"os"
	// "golang.org/x/net/websocket"
	// "fmt"
	"log"
	// "encoding/json"
	// "os"
	"os/user"
	"github.com/gorilla/websocket"
)

// websocket客户端
type WebsocketClient struct {
	Url string
	OnMessage chan string
}

type Message struct {
	Type       string `json:"type"`
	ClientName string `json:"client_name"`
	Data       string `json:"data"`
}

func (client *WebsocketClient) Start() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	Info("Connecting to %s", client.Url)

	c, _, err := websocket.DefaultDialer.Dial(client.Url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// // login
	current, _ := user.Current()
	var message = Message{"login", current.Username, ""}
	// // var str, _ = json.Marshal(message)
	// fmt.Println(message)
	c.WriteJSON(message)
	// c.WriteMessage(websocket.TextMessage, []byte("{\"type\":\"login\",\"client_name\":\"" + current.Username + "\"}"))

	done := make(chan struct{})
	chanMsg := make(chan string)
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			chanMsg <- string(message)
			// j := make(map[string]interface{})
			// json.Unmarshal(message, &j)

		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close: ", err)
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
