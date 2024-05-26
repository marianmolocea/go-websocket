package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
	upgrader websocket.Upgrader
}

func homePage(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "Home")
	log.Println("GET Home")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Websocket")
}

func (wsh WebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s wehen upgrading connection to websocket", err)
		return
	}

	defer func() {
		log.Println("closing connection")
		c.Close()
	}()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			return
		}
		if mt == websocket.BinaryMessage {
			err = c.WriteMessage(websocket.TextMessage, []byte("server doesn't support binary messages"))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
			}
			return
		}
		log.Printf("Receive message %s", string(message))

		if strings.Trim(string(message), "\n") == "stop" {
			log.Printf("received stop message")
			err = c.WriteMessage(websocket.TextMessage, []byte("Disconnected"))
			if err != nil {
				defer c.Close()
			}

			return
		}

		if strings.Trim(string(message), "\n") != "start" {
			err = c.WriteMessage(websocket.TextMessage, []byte("You did not say the magic word"))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
				return
			}
			continue
		}
		log.Println("start responding to client...")
	}
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	webSoocketHandler := WebsocketHandler{
		upgrader: websocket.Upgrader{},
	}
	http.Handle("/ws", webSoocketHandler)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
