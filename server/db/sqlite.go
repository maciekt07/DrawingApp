package db

import (
	"database/sql"

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
        x FLOAT,
        y FLOAT,
        color TEXT
    );`
    _, err = DB.Exec(createTable)
    if err != nil {
        panic(err)
    }
}
