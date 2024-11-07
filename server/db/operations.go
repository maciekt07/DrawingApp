package db

import (
	"encoding/json"
	"log"
	"server/models"
)

//  clears all drawings from the database
func ClearDrawingsInDB() {
	_, err := DB.Exec("DELETE FROM drawings")
	if err != nil {
		log.Println("Error clearing drawings:", err)
	}
}

// saves a drawing to the database
func SaveDrawingToDB(msg models.DrawingMessage) {
	stmt, err := DB.Prepare("INSERT INTO drawings (path, color) VALUES (?, ?)")
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
