package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/db"
	"server/models"
)


func SaveDrawingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var drawing models.Drawing
	if err := json.NewDecoder(r.Body).Decode(&drawing); err != nil {
		http.Error(w, "Invalid drawing data", http.StatusBadRequest)
		return
	}

	// Save drawing to database
	stmt, err := db.DB.Prepare("INSERT INTO drawings (x, y, color) VALUES (?, ?, ?)")
    fmt.Println(stmt)
	if err != nil {
		log.Println("Error preparing statement:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	if _, err := stmt.Exec(drawing.X, drawing.Y, drawing.Color); err != nil {
		log.Println("Error executing statement:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetDrawingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.DB.Query("SELECT id, x, y, color FROM drawings")
	if err != nil {
		log.Println("Error querying database:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var drawings []models.Drawing
	for rows.Next() {
		var drawing models.Drawing
		if err := rows.Scan(&drawing.ID, &drawing.X, &drawing.Y, &drawing.Color); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		drawings = append(drawings, drawing)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drawings)
}
