package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

// DB wraps the standard sql.DB to provide additional functionality.
type DB struct {
	*sql.DB
}

// Initialize sets up a new database connection and initializes the required tables.
func Initialize(databasePath string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := createTables(sqlDB); err != nil {
		return nil, err
	}

	return &DB{sqlDB}, nil
}

// createTables creates the necessary tables in the database if they don't already exist.
func createTables(db *sql.DB) error {
	// Create sessions table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			session_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT DEFAULT 'Session'
		);`); err != nil {
		return fmt.Errorf("failed to create sessions table: %w", err)
	}

	// Create history table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id INTEGER,
			user_input TEXT,
			ai_response TEXT,
			FOREIGN KEY(session_id) REFERENCES sessions(session_id)
		);`); err != nil {
		return fmt.Errorf("failed to create history table: %w", err)
	}

	return nil
}

func (db *DB) SQLDB() *sql.DB {
	return db.DB
}
