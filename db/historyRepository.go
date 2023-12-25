package db

import (
	"database/sql"
)

type HistoryRepository struct {
	db *sql.DB
}

func NewHistoryRepository(db *sql.DB) *HistoryRepository {
	return &HistoryRepository{db: db}
}

// GetSessionHistory retrieves the history for a given session.
func (repo *HistoryRepository) GetSessionHistory(sessionID string) ([]map[string]string, error) {
	rows, err := repo.db.Query("SELECT user_input, ai_response FROM history WHERE session_id = ? ORDER BY id ASC", sessionID)
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

// InsertChatHistory adds a new entry to the chat history.
func (repo *HistoryRepository) InsertChatHistory(sessionID, userInput, aiResponse string) error {
	_, err := repo.db.Exec("INSERT INTO history (session_id, user_input, ai_response) VALUES (?, ?, ?)", sessionID, userInput, aiResponse)
	return err
}
