package services

import (
	"GPTChat/db"
)

type ChatService struct {
	HistoryRepo *db.HistoryRepository
}

func NewChatService(historyRepo *db.HistoryRepository) *ChatService {
	return &ChatService{HistoryRepo: historyRepo}
}

// GetSessionHistory retrieves the chat history for a given session.
func (s *ChatService) GetSessionHistory(sessionID string) ([]map[string]string, error) {
	return s.HistoryRepo.GetSessionHistory(sessionID)
}

// SaveChatHistory adds a new entry to the chat history.
func (s *ChatService) SaveChatHistory(sessionID, userInput, aiResponse string) error {
	return s.HistoryRepo.InsertChatHistory(sessionID, userInput, aiResponse)
}
