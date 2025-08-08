package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var GDB *sql.DB

func Init() {
	var lErr error

	GDB, lErr = sql.Open("sqlite", "./data/mydb.db")
	if lErr != nil {
		log.Fatal("Failed to connect to database Error:DBI01", lErr)
	} else {
		log.Println("Connected succesfullly")
	}
	createTable := `
	CREATE TABLE IF NOT EXISTS audio_chunks (
    chunk_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT,
    session_id TEXT,
    timestamp TEXT,
    file_name TEXT,
    content_type TEXT,
    audio_data BLOB,
    duration TEXT,
    size INTEGER,
    transcript TEXT,
    checksum TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`
	_, lErr = GDB.Exec(createTable)
	if lErr != nil {
		log.Fatal("Error:DBI01", lErr)
	}

}
