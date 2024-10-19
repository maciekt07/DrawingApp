package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./drawings.db")
	if err != nil {
		panic(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS drawings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT,
		color TEXT
	);`
	_, err = DB.Exec(createTable)
	if err != nil {
		panic(err)
	}
    log.Println("Database connected successfully.")
}
