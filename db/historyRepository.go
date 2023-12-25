package db

import (
	"database/sql"
	"fmt"
)

// HistoryRepository handles operations on chat history data in the database.
type HistoryRepository struct {
	db *sql.DB
}

// NewHistoryRepository creates a new instance of HistoryRepository.
func NewHistoryRepository(db *sql.DB) *HistoryRepository {
	return &HistoryRepository{db: db}
}

// GetSessionHistory retrieves the chat history for a given session.
func (repo *HistoryRepository) GetSessionHistory(sessionID string) ([]map[string]string, error) {
	rows, err := repo.db.Query("SELECT user_input, ai_response FROM history WHERE session_id = ? ORDER BY id ASC", sessionID)
	if err != nil {
		return nil, fmt.Errorf("error querying session history for session_id %s: %w", sessionID, err)
	}
	defer rows.Close()

	var history []map[string]string
	for rows.Next() {
		var userInput, aiResponse string
		if err := rows.Scan(&userInput, &aiResponse); err != nil {
			return nil, fmt.Errorf("error scanning history row: %w", err)
		}

		if userInput != "" {
			userEntry := map[string]string{"role": "user", "content": userInput}
			history = append(history, userEntry)
		}

		if aiResponse != "" {
			aiEntry := map[string]string{"role": "assistant", "content": aiResponse}
			history = append(history, aiEntry)
		}
	}
	return history, nil
}

// InsertChatHistory adds a new entry to the chat history in the database.
func (repo *HistoryRepository) InsertChatHistory(sessionID, userInput, aiResponse string) error {
	if _, err := repo.db.Exec("INSERT INTO history (session_id, user_input, ai_response) VALUES (?, ?, ?)", sessionID, userInput, aiResponse); err != nil {
		return fmt.Errorf("error inserting chat history for session_id %s: %w", sessionID, err)
	}
	return nil
}
