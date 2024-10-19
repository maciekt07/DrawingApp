package main

import (
	"log"
	"net/http"
	"server/db"
	"server/handlers"

	"github.com/gorilla/mux"
)

// CORS middleware
func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        // Handle preflight requests
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK) // Respond to preflight
            return
        }

        next.ServeHTTP(w, r)
    })
}


func main() {
    // Initialize the database
    db.InitDB()

    // Set up the router
    r := mux.NewRouter()

    // WebSocket route
    r.HandleFunc("/ws", handlers.WebSocketHandler)

    // Drawings route
    r.HandleFunc("/drawings", handlers.SaveDrawingHandler).Methods("POST")
    r.HandleFunc("/drawings", handlers.GetDrawingsHandler).Methods("GET")
    

    // Apply CORS middleware
    http.Handle("/", enableCORS(r))

    // Start the server
    log.Println("Server started on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
