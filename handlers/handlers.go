package handlers

import (
	"GPTChat/db"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Handler struct {
	DB *db.DB
}

func NewHandler(db *db.DB) *Handler {
	return &Handler{DB: db}
}

// Utility to send JSON response
func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Utility to send error response
func sendErrorResponse(w http.ResponseWriter, errorMessage string, statusCode int) {
	sendJSONResponse(w, map[string]string{"error": errorMessage}, statusCode)
}

func (h *Handler) LogHandler(w http.ResponseWriter, r *http.Request) {
	var logObject struct {
		Message string `json:"message"`
		IsError bool   `json:"isError"`
	}
	if err := json.NewDecoder(r.Body).Decode(&logObject); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	logMessage := logObject.Message

	file, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(logMessage + "\n"); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, map[string]string{"status": "logged"}, http.StatusOK)
}

func (h *Handler) QueryHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		sendErrorResponse(w, "session_id is required", http.StatusBadRequest)
		return
	}
	userQuery := r.URL.Query().Get("query")

	// Retrieve and log session history
	history, err := h.DB.GetSessionHistory(sessionID)
	if err != nil {
		log.Printf("Error retrieving session history for session %s: %v", sessionID, err)
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Append the current user query
	history = append(history, map[string]string{"role": "user", "content": userQuery})
	historyJSON, err := json.Marshal(history)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Run the Python script and get the output
	cmd := exec.Command("python", "main.py")
	cmd.Stdin = strings.NewReader(string(historyJSON))
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error running AI model for session %s: %v", sessionID, err)
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert the AI response into the database and log the operation
	if err := h.DB.InsertChatHistory(sessionID, userQuery, string(output)); err != nil {
		log.Printf("Error saving chat history for session %s: %v", sessionID, err)
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the plain output directly to the client
	w.Header().Set("Content-Type", "text/plain")
	w.Write(output)
}

func (h *Handler) NewSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionID, err := h.DB.InsertNewSession()
	if err != nil {
		log.Printf("Error creating new session: %v", err)
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, map[string]int64{"session_id": sessionID}, http.StatusOK)
}

func (h *Handler) GetSessionsHandler(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.DB.GetAllSessions()
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, sessions, http.StatusOK)
}

func (h *Handler) DeleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	err := h.DB.DeleteSession(sessionID)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
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
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.DB.RenameSession(request.SessionID, request.NewName)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, map[string]string{"status": "renamed"}, http.StatusOK)
}

func (h *Handler) GetSessionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		sendErrorResponse(w, "session_id is required", http.StatusBadRequest)
		return
	}
	history, err := h.DB.GetSessionHistory(sessionID)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, history, http.StatusOK)
}
