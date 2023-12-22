package main

import (
	"GPTChat/db"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	db.Initialize("./db/conversation.db")

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/log", logHandler)
	http.HandleFunc("/query", queryHandler)
	http.HandleFunc("/new-session", newSessionHandler)
	http.HandleFunc("/get-sessions", getSessionsHandler)
	http.HandleFunc("/get-session-history", getSessionHistoryHandler)
	http.HandleFunc("/delete-session", deleteSessionHandler)
	http.HandleFunc("/rename-session", renameSessionHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	var logObject struct {
		Message string `json:"message"`
		IsError bool   `json:"isError"`
	}
	if err := json.NewDecoder(r.Body).Decode(&logObject); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logMessage := logObject.Message

	file, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(logMessage + "\n"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "session_id is required", http.StatusBadRequest)
		return
	}
	userQuery := r.URL.Query().Get("query")

	// Retrieve and log session history
	history, err := db.GetSessionHistory(sessionID)
	if err != nil {
		log.Printf("Error retrieving session history for session %s: %v", sessionID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Append the current user query
	history = append(history, map[string]string{"role": "user", "content": userQuery})
	historyJSON, _ := json.Marshal(history)

	// Run the Python script
	cmd := exec.Command("python", "main.py")
	cmd.Stdin = strings.NewReader(string(historyJSON))
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error running AI model for session %s: %v", sessionID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert the AI response into the database and log the operation
	err = db.InsertChatHistory(sessionID, userQuery, string(output))
	if err != nil {
		log.Printf("Error saving chat history for session %s: %v", sessionID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(output)
}

func newSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionID, err := db.InsertNewSession()
	if err != nil {
		log.Printf("Error creating new session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int64{"session_id": sessionID})
}

func getSessionsHandler(w http.ResponseWriter, r *http.Request) {
	sessions, err := db.GetAllSessions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(sessions)
}

func deleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	err := db.DeleteSession(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func renameSessionHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SessionID string `json:"sessionId"`
		NewName   string `json:"newName"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.RenameSession(request.SessionID, request.NewName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getSessionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "session_id is required", http.StatusBadRequest)
		return
	}
	history, err := db.GetSessionHistory(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(history)
}
