package handlers

import (
	"GPTChat/errors"
	"GPTChat/services"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Handler encapsulates the dependencies for handling HTTP requests.
type Handler struct {
	chatService    *services.ChatService
	sessionService *services.SessionService
}

// NewHandler constructs a new Handler with the given repositories.
func NewHandler(chatService *services.ChatService, sessionService *services.SessionService) *Handler {
	return &Handler{
		chatService:    chatService,
		sessionService: sessionService,
	}
}

// sendJSONResponse is a utility function to send a JSON response.
func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

// LogHandler handles logging messages to a file.
func (h *Handler) LogHandler(w http.ResponseWriter, r *http.Request) {
	var logObject struct {
		Message string `json:"message"`
		IsError bool   `json:"isError"`
	}

	if err := json.NewDecoder(r.Body).Decode(&logObject); err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusBadRequest, "Error decoding JSON", "An unexpected error has occurred."))
		return
	}

	if err := appendToFile("application.log", logObject.Message+"\n"); err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, err.Error(), "An unexpected error has occurred."))
		return
	}
	sendJSONResponse(w, map[string]string{"status": "logged"}, http.StatusOK)
}

// appendToFile writes a message to the specified file.
func appendToFile(filename, message string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(message); err != nil {
		return err
	}
	return nil
}

// QueryHandler handles the querying process.
func (h *Handler) QueryHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		errors.HandleError(w, errors.NewCustomError(http.StatusBadRequest, "session_id is required", "An unexpected error has occurred."))
		return
	}
	userQuery := r.URL.Query().Get("query")

	history, err := h.chatService.GetSessionHistory(sessionID)
	if err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error retrieving session history", "An unexpected error has occurred."))
		return
	}

	history = append(history, map[string]string{"role": "user", "content": userQuery})
	historyJSON, err := json.Marshal(history)
	if err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error marshalling session history", "An unexpected error has occurred."))
		return
	}

	cmd := exec.Command("python", "main.py")
	cmd.Stdin = strings.NewReader(string(historyJSON))
	output, err := cmd.Output()
	if err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error running AI model", "An unexpected error has occurred."))
		return
	}

	if err := h.chatService.SaveChatHistory(sessionID, userQuery, string(output)); err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error saving chat history", "An unexpected error has occurred."))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(output)
}

// NewSessionHandler handles the creation of a new session.
func (h *Handler) NewSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionID, err := h.sessionService.CreateSession()
	if err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error creating new session", "An unexpected error has occurred."))
		return
	}

	sendJSONResponse(w, map[string]int64{"session_id": sessionID}, http.StatusOK)
}

// GetSessionsHandler handles retrieving all sessions.
func (h *Handler) GetSessionsHandler(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.sessionService.GetAllSessions()
	if err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error retrieving sessions", "An unexpected error has occurred."))
		return
	}

	sendJSONResponse(w, sessions, http.StatusOK)
}

// DeleteSessionHandler handles deleting a session.
func (h *Handler) DeleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if err := h.sessionService.DeleteSession(sessionID); err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error deleting session", "An unexpected error has occurred."))
		return
	}
	sendJSONResponse(w, map[string]string{"status": "deleted"}, http.StatusOK)
}

// RenameSessionHandler handles renaming a session.
func (h *Handler) RenameSessionHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SessionID string `json:"sessionId"`
		NewName   string `json:"newName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusBadRequest, "Error decoding JSON", "An unexpected error has occurred."))
		return
	}

	if err := h.sessionService.RenameSession(request.SessionID, request.NewName); err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error renaming session", "An unexpected error has occurred."))
		return
	}
	sendJSONResponse(w, map[string]string{"status": "renamed"}, http.StatusOK)
}

// GetSessionHistoryHandler handles retrieving session history.
func (h *Handler) GetSessionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		errors.HandleError(w, errors.NewCustomError(http.StatusBadRequest, "session_id is required", "An unexpected error has occurred."))
		return
	}

	history, err := h.chatService.GetSessionHistory(sessionID)
	if err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error retrieving session history", "An unexpected error has occurred."))
		return
	}
	sendJSONResponse(w, history, http.StatusOK)
}
