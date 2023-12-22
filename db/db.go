package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type DB struct {
	*sql.DB
}

func Initialize(databasePath string) (*DB, error) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, err // return the error instead of exiting
	}

	// Create sessions table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			session_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT DEFAULT 'Session'
		);
	`)
	if err != nil {
		return nil, err // return the error instead of exiting
	}

	// Create history table if it doesn't exist
	_, err = db.Exec(`
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

	return &DB{db}, nil // return your DB instance
}

// GetSessionHistory retrieves the history for a given session.
func (db *DB) GetSessionHistory(sessionID string) ([]map[string]string, error) {
	rows, err := db.DB.Query("SELECT user_input, ai_response FROM history WHERE session_id = ? ORDER BY id ASC", sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []map[string]string
	for rows.Next() {
		var userInput, aiResponse string
		err := rows.Scan(&userInput, &aiResponse)
		if err != nil {
			return nil, err
		}

		// Append user message to history if not empty
		if userInput != "" {
			userEntry := map[string]string{
				"role":    "user",
				"content": userInput,
			}
			history = append(history, userEntry)
		}

		// Append AI response to history if not empty
		if aiResponse != "" {
			aiEntry := map[string]string{
				"role":    "assistant",
				"content": aiResponse,
			}
			history = append(history, aiEntry)
		}
	}
	return history, nil
}

func (db *DB) InsertNewSession() (int64, error) {
	result, err := db.DB.Exec("INSERT INTO sessions DEFAULT VALUES")
	if err != nil {
		return 0, err
	}
	sessionID, err := result.LastInsertId()
	return sessionID, err
}

func (db *DB) DeleteSession(sessionID string) error {
	_, err := db.DB.Exec("DELETE FROM history WHERE session_id = ?", sessionID)
	if err != nil {
		return err
	}
	_, err = db.DB.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	return err
}

func (db *DB) RenameSession(sessionID, newName string) error {
	_, err := db.DB.Exec("UPDATE sessions SET name = ? WHERE session_id = ?", newName, sessionID)
	return err
}

// InsertChatHistory adds a new entry to the chat history.
func (db *DB) InsertChatHistory(sessionID, userInput, aiResponse string) error {
	_, err := db.DB.Exec("INSERT INTO history (session_id, user_input, ai_response) VALUES (?, ?, ?)", sessionID, userInput, aiResponse)
	return err
}

// GetAllSessions retrieves all sessions from the database.
func (db *DB) GetAllSessions() ([]map[string]interface{}, error) {
	rows, err := db.DB.Query("SELECT session_id, name FROM sessions ORDER BY session_id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []map[string]interface{}
	for rows.Next() {
		var sessionID int64
		var name string
		err := rows.Scan(&sessionID, &name)
		if err != nil {
			return nil, err
		}
		session := map[string]interface{}{"session_id": sessionID, "name": name}
		sessions = append(sessions, session)
	}
	return sessions, nil
}
