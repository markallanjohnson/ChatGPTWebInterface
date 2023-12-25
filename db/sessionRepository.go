// Package db provides data access and manipulation functions for session and history data.
package db

import (
	"database/sql"
	"fmt"
)

// SessionRepository handles operations on session data in the database.
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new instance of SessionRepository.
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// InsertNewSession inserts a new session into the database and returns its ID.
func (repo *SessionRepository) InsertNewSession() (int64, error) {
	result, err := repo.db.Exec("INSERT INTO sessions DEFAULT VALUES")
	if err != nil {
		return 0, fmt.Errorf("error inserting new session: %w", err)
	}
	return result.LastInsertId()
}

// DeleteSession removes a session and its related data from the database.
func (repo *SessionRepository) DeleteSession(sessionID string) error {
	if _, err := repo.db.Exec("DELETE FROM history WHERE session_id = ?", sessionID); err != nil {
		return fmt.Errorf("error deleting session history for session_id %s: %w", sessionID, err)
	}
	if _, err := repo.db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID); err != nil {
		return fmt.Errorf("error deleting session for session_id %s: %w", sessionID, err)
	}
	return nil
}

// RenameSession updates the name of a session in the database.
func (repo *SessionRepository) RenameSession(sessionID, newName string) error {
	if _, err := repo.db.Exec("UPDATE sessions SET name = ? WHERE session_id = ?", newName, sessionID); err != nil {
		return fmt.Errorf("error renaming session for session_id %s: %w", sessionID, err)
	}
	return nil
}

// GetAllSessions retrieves all session data from the database.
func (repo *SessionRepository) GetAllSessions() ([]map[string]interface{}, error) {
	rows, err := repo.db.Query("SELECT session_id, name FROM sessions ORDER BY session_id DESC")
	if err != nil {
		return nil, fmt.Errorf("error retrieving all sessions: %w", err)
	}
	defer rows.Close()

	var sessions []map[string]interface{}
	for rows.Next() {
		var sessionID int64
		var name string
		if err := rows.Scan(&sessionID, &name); err != nil {
			return nil, fmt.Errorf("error scanning session data: %w", err)
		}
		sessions = append(sessions, map[string]interface{}{"session_id": sessionID, "name": name})
	}
	return sessions, nil
}
