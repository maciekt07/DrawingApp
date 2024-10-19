package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)

type DrawingMessage struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Color string  `json:"color"`
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	defer delete(clients, conn)

	log.Println("Client connected:", conn.RemoteAddr())

	for {
		var msg DrawingMessage
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Error while reading message:", err)
			break
		}
		broadcast(msg)
	}
}

func broadcast(msg DrawingMessage) {
	for client := range clients {
		if err := client.WriteJSON(msg); err != nil {
			log.Printf("Error while broadcasting message to client %v: %v", client.RemoteAddr(), err)
			client.Close()
			delete(clients, client)
		}
	}
}
