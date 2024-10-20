package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/db"

	"github.com/gorilla/websocket"
)

// upgrade the http connection to a websocket connection
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}
// active websocket connections
var clients = make(map[*websocket.Conn]bool)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type DrawingMessage struct {
	Path  []Point `json:"path"`
	Color string  `json:"color"`
}

// handles incoming WebSocket connections
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
    // listen for incoming messages from the client
    for {
        var msg DrawingMessage
        if err := conn.ReadJSON(&msg); err != nil {
            log.Println("Error while reading message:", err)
            break
        }

        // Save the drawing message directly to the database
        saveDrawingToDB(msg)

        broadcast(msg)
    }
}

// TODO: move this somewhere else
func saveDrawingToDB(msg DrawingMessage) {
    stmt, err := db.DB.Prepare("INSERT INTO drawings (path, color) VALUES (?, ?)")
    if err != nil {
        log.Println("Error preparing statement:", err)
        return
    }
    defer stmt.Close()

    pathJSON, err := json.Marshal(msg.Path)
    if err != nil {
        log.Println("Error marshalling path:", err)
        return
    }

    if _, err := stmt.Exec(pathJSON, msg.Color); err != nil {
        log.Println("Error executing statement:", err)
    }
}

// sends the drawing message to all connected clients
func broadcast(msg DrawingMessage) {
	for client := range clients {
		if err := client.WriteJSON(msg); err != nil {
			log.Printf("Error while broadcasting message to client %v: %v", client.RemoteAddr(), err)
			client.Close()
			delete(clients, client)
		}
	}
}