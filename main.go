package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./conversation.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			session_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT DEFAULT 'New Session'
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id INTEGER,
			user_input TEXT,
			ai_response TEXT,
			FOREIGN KEY(session_id) REFERENCES sessions(session_id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func getSessionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")

	rows, err := db.Query("SELECT user_input, ai_response FROM history WHERE session_id = ?", sessionID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var conversationHistory []map[string]string
	for rows.Next() {
		var userInput, aiResponse string
		if err := rows.Scan(&userInput, &aiResponse); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if userInput != "" {
			conversationHistory = append(conversationHistory, map[string]string{"role": "user", "content": userInput})
		}
		if aiResponse != "" {
			conversationHistory = append(conversationHistory, map[string]string{"role": "assistant", "content": aiResponse})
		}
	}

	json.NewEncoder(w).Encode(conversationHistory)
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

func main() {
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

func queryHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	userQuery := r.URL.Query().Get("query")

	// Retrieve conversation history from the database for the given session
	rows, err := db.Query("SELECT user_input, ai_response FROM history WHERE session_id = ?", sessionID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var conversationHistory []map[string]string
	for rows.Next() {
		var userInput, aiResponse string
		if err := rows.Scan(&userInput, &aiResponse); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if userInput != "" {
			conversationHistory = append(conversationHistory, map[string]string{"role": "user", "content": userInput})
		}
		if aiResponse != "" {
			conversationHistory = append(conversationHistory, map[string]string{"role": "assistant", "content": aiResponse})
		}
	}

	// Append the current user query
	conversationHistory = append(conversationHistory, map[string]string{"role": "user", "content": userQuery})

	// Convert conversation history to JSON for the Python script
	historyJSON, _ := json.Marshal(conversationHistory)

	// Run the Python script
	cmd := exec.Command("python", "main.py")
	cmd.Stdin = strings.NewReader(string(historyJSON))
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Append the AI response to the database
	_, err = db.Exec("INSERT INTO history (session_id, user_input, ai_response) VALUES (?, ?, ?)", sessionID, userQuery, string(output))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Send the Python script's output back to the client
	w.Write(output)
}

func newSessionHandler(w http.ResponseWriter, r *http.Request) {
	result, err := db.Exec("INSERT INTO sessions DEFAULT VALUES")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	sessionId, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]int64{"session_id": sessionId})
}

func getSessionsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT session_id, name FROM sessions ORDER BY session_id DESC")
	if err != nil {
		log.Printf("SQL error in getSessionsHandler: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var sessions []map[string]interface{}
	for rows.Next() {
		var sessionID int64
		var name string
		if err := rows.Scan(&sessionID, &name); err != nil {
			log.Printf("SQL scan error in getSessionsHandler: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}
		sessions = append(sessions, map[string]interface{}{"session_id": sessionID, "name": name})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sessions); err != nil {
		log.Printf("JSON encoding error in getSessionsHandler: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func deleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	_, err := db.Exec("DELETE FROM history WHERE session_id = ?", sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
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
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Use SQL parameters to prevent SQL injection.
	_, err := db.Exec("UPDATE sessions SET name = ? WHERE session_id = ?", request.NewName, request.SessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
