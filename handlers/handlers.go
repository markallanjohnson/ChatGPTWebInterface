package handlers

import (
	"GPTChat/db"
	"GPTChat/errors"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Handler struct {
	sessionRepo *db.SessionRepository
	historyRepo *db.HistoryRepository
}

func NewHandler(sessionRepo *db.SessionRepository, historyRepo *db.HistoryRepository) *Handler {
	return &Handler{
		sessionRepo: sessionRepo,
		historyRepo: historyRepo,
	}
}

// Utility to send JSON response
func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) LogHandler(w http.ResponseWriter, r *http.Request) {
	var logObject struct {
		Message string `json:"message"`
		IsError bool   `json:"isError"`
	}
	if err := json.NewDecoder(r.Body).Decode(&logObject); err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusBadRequest, "Error decoding JSON", "An unexpected error has occurred."))
		return
	}
	logMessage := logObject.Message

	file, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error opening log file", "An unexpected error has occurred."))
		return
	}
	defer file.Close()

	if _, err := file.WriteString(logMessage + "\n"); err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error writing to log file", "An unexpected error has occurred."))
		return
	}
	sendJSONResponse(w, map[string]string{"status": "logged"}, http.StatusOK)
}

func (h *Handler) QueryHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		errors.HandleError(w, errors.NewCustomError(http.StatusBadRequest, "session_id is required", "An unexpected error has occurred."))
		return
	}
	userQuery := r.URL.Query().Get("query")

	// Retrieve and log session history
	history, err := h.historyRepo.GetSessionHistory(sessionID)
	if err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error retrieving session history", "An unexpected error has occurred."))
		return
	}

	// Append the current user query
	history = append(history, map[string]string{"role": "user", "content": userQuery})
	historyJSON, err := json.Marshal(history)
	if err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error marshalling session history", "An unexpected error has occurred."))
		return
	}

	// Run the Python script and get the output
	cmd := exec.Command("python", "main.py")
	cmd.Stdin = strings.NewReader(string(historyJSON))
	output, err := cmd.Output()
	if err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error running AI model", "An unexpected error has occurred."))
		return
	}

	// Insert the AI response into the database and log the operation
	if err := h.historyRepo.InsertChatHistory(sessionID, userQuery, string(output)); err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error saving chat history", "An unexpected error has occurred."))
		return
	}

	// Send the plain output directly to the client
	w.Header().Set("Content-Type", "text/plain")
	w.Write(output)
}

func (h *Handler) NewSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionID, err := h.sessionRepo.InsertNewSession()
	if err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error creating new session", "An unexpected error has occurred."))
		return
	}

	sendJSONResponse(w, map[string]int64{"session_id": sessionID}, http.StatusOK)
}

func (h *Handler) GetSessionsHandler(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.sessionRepo.GetAllSessions()
	if err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error retrieving sessions", "An unexpected error has occurred."))
		return
	}

	sendJSONResponse(w, sessions, http.StatusOK)
}

func (h *Handler) DeleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	err := h.sessionRepo.DeleteSession(sessionID)
	if err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error deleting session", "An unexpected error has occurred."))
		return
	}
	sendJSONResponse(w, map[string]string{"status": "deleted"}, http.StatusOK)
}

func (h *Handler) RenameSessionHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SessionID string `json:"sessionId"`
		NewName   string `json:"newName"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusBadRequest, "Error decoding JSON", "An unexpected error has occurred."))
		return
	}

	err = h.sessionRepo.RenameSession(request.SessionID, request.NewName)
	if err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error renaming session", "An unexpected error has occurred."))
		return
	}
	sendJSONResponse(w, map[string]string{"status": "renamed"}, http.StatusOK)
}

func (h *Handler) GetSessionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		errors.HandleError(w, errors.NewCustomError(http.StatusBadRequest, "session_id is required", "An unexpected error has occurred."))
		return
	}
	history, err := h.historyRepo.GetSessionHistory(sessionID)
	if err != nil {
		errors.HandleError(w, errors.NewCustomError(http.StatusInternalServerError, "Error retrieving session history", "An unexpected error has occurred."))
		return
	}
	sendJSONResponse(w, history, http.StatusOK)
}
