package handlers

import (
	"log"
	"net/http"
	"server/db"
	"server/models"

	"github.com/gorilla/websocket"
)

// upgrade the http connection to a WebSocket connection
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

// active WebSocket connections
var clients = make(map[*websocket.Conn]bool)

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
        var msg models.DrawingMessage
        if err := conn.ReadJSON(&msg); err != nil {
            log.Println("Error while reading message:", err)
            break
        }

        if msg.Type == "clear" {
            db.ClearDrawingsInDB()   //  clear drawings from the database
            broadcastClearEvent()    //  clear event to all clients
        } else {
            // save the drawing message directly to the database
            db.SaveDrawingToDB(msg)
            broadcast(msg)
        }
    }
}

// broadcasts a clear event to all connected clients
func broadcastClearEvent() {
    msg := models.DrawingMessage{Type: "clear"}
    for client := range clients {
        if err := client.WriteJSON(msg); err != nil {
            log.Printf("Error while broadcasting clear event to client %v: %v", client.RemoteAddr(), err)
            client.Close()
            delete(clients, client)
        }
    }
}

// sends the drawing message to all connected clients
func broadcast(msg models.DrawingMessage) {
	for client := range clients {
		if err := client.WriteJSON(msg); err != nil {
			log.Printf("Error while broadcasting message to client %v: %v", client.RemoteAddr(), err)
			client.Close()
			delete(clients, client)
		}
	}
}

// TODO: add acive users count
// func activeUsersHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	var mu sync.Mutex
// 	mu.Lock()
// 	count := len(clients)
// 	mu.Unlock()

// 	// Respond with the count of active users
// 	response := map[string]int{"active_users": count}
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(response)
// }
