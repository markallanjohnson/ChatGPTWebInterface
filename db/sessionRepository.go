package db

import (
	"database/sql"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (repo *SessionRepository) InsertNewSession() (int64, error) {
	result, err := repo.db.Exec("INSERT INTO sessions DEFAULT VALUES")
	if err != nil {
		return 0, err
	}
	sessionID, err := result.LastInsertId()
	return sessionID, err
}

func (repo *SessionRepository) DeleteSession(sessionID string) error {
	_, err := repo.db.Exec("DELETE FROM history WHERE session_id = ?", sessionID)
	if err != nil {
		return err
	}
	_, err = repo.db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	return err
}

func (repo *SessionRepository) RenameSession(sessionID, newName string) error {
	_, err := repo.db.Exec("UPDATE sessions SET name = ? WHERE session_id = ?", newName, sessionID)
	return err
}

// GetAllSessions retrieves all sessions from the database.
func (repo *SessionRepository) GetAllSessions() ([]map[string]interface{}, error) {
	rows, err := repo.db.Query("SELECT session_id, name FROM sessions ORDER BY session_id DESC")
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
