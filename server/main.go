package main

import (
	"encoding/json"
	"log"
	"net/http"
	"server/db"
	"server/handlers"
	"server/middleware"

	"github.com/gorilla/mux"
)

func main() {
    // Initialize the database
    db.InitDB()

    // Set up the router
    r := mux.NewRouter()

    // WebSocket route
    r.HandleFunc("/ws", handlers.WebSocketHandler)

    // Drawings route
    r.HandleFunc("/drawings", handlers.SaveDrawingHandler).Methods("POST") // this is not being used
    r.HandleFunc("/drawings", handlers.GetDrawingsHandler).Methods("GET")
    
    http.HandleFunc("/health", healthCheckHandler)

    // Apply CORS middleware
    http.Handle("/", middleware.EnableCORS(r))

    // Start the server
    log.Println("Server started on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Set the header for JSON response
	w.Header().Set("Content-Type", "application/json")
	
	// Respond with a simple JSON object
	response := map[string]string{"status": "ok"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
