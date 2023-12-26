package services

import (
	"GPTChat/db"
)

type SessionService struct {
	SessionRepo *db.SessionRepository
}

func NewSessionService(sessionRepo *db.SessionRepository) *SessionService {
	return &SessionService{SessionRepo: sessionRepo}
}

// CreateSession creates a new session in the database.
func (s *SessionService) CreateSession() (int64, error) {
	return s.SessionRepo.InsertNewSession()
}

// GetAllSessions retrieves all sessions from the database.
func (s *SessionService) GetAllSessions() ([]map[string]interface{}, error) {
	return s.SessionRepo.GetAllSessions()
}

// DeleteSession removes a session from the database.
func (s *SessionService) DeleteSession(sessionID string) error {
	return s.SessionRepo.DeleteSession(sessionID)
}

// RenameSession changes the name of an existing session.
func (s *SessionService) RenameSession(sessionID, newName string) error {
	return s.SessionRepo.RenameSession(sessionID, newName)
}
