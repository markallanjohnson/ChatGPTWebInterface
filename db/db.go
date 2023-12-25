package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type DB struct {
	*sql.DB
}

func Initialize(databasePath string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, err // return the error instead of exiting
	}

	// Create sessions table if it doesn't exist
	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			session_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT DEFAULT 'Session'
		);
	`)
	if err != nil {
		return nil, err // return the error instead of exiting
	}

	// Create history table if it doesn't exist
	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id INTEGER,
			user_input TEXT,
			ai_response TEXT,
			FOREIGN KEY(session_id) REFERENCES sessions(session_id)
		);
	`)
	if err != nil {
		return nil, err // return the error instead of exiting
	}

	return &DB{sqlDB}, nil // return your DB instance
}

func (db *DB) SQLDB() *sql.DB {
	return db.DB
}
