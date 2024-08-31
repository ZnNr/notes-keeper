package repository

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateTables(db *sql.DB) {
	userTable := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );`

	noteTable := `CREATE TABLE IF NOT EXISTS notes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        text TEXT NOT NULL,
        mistakes TEXT,
        FOREIGN KEY(user_id) REFERENCES users(id)
    );`

	if _, err := db.Exec(userTable); err != nil {
		log.Fatalf("Failed to create users table: %s", err)
	}

	if _, err := db.Exec(noteTable); err != nil {
		log.Fatalf("Failed to create notes table: %s", err)
	}

	fmt.Println("Tables created successfully")
}
