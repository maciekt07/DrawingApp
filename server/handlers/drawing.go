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
    fmt.Println("SaveDrawingHandler")
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var drawing models.Drawing
    if err := json.NewDecoder(r.Body).Decode(&drawing); err != nil {
        http.Error(w, "Invalid drawing data", http.StatusBadRequest)
        return
    }

    stmt, err := db.DB.Prepare("INSERT INTO drawings (path, color) VALUES (?, ?)")
    if err != nil {
        log.Println("Error preparing statement:", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    pathJSON, _ := json.Marshal(drawing.Path)
    if _, err := stmt.Exec(pathJSON, drawing.Color); err != nil {
        log.Println("Error executing statement:", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

// retrieves all drawings from the database
func GetDrawingsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, path, color FROM drawings")
	if err != nil {
		http.Error(w, "Error retrieving drawings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var drawings []models.Drawing
	for rows.Next() {
		var drawing models.Drawing
		var pathJSON string

		if err := rows.Scan(&drawing.ID, &pathJSON, &drawing.Color); err != nil {
			http.Error(w, "Error scanning drawing", http.StatusInternalServerError)
			return
		}

		if err := json.Unmarshal([]byte(pathJSON), &drawing.Path); err != nil {
			http.Error(w, "Error unmarshalling drawing path", http.StatusInternalServerError)
			return
		}

		drawings = append(drawings, drawing)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drawings)
}