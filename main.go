package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
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
			session_id INTEGER PRIMARY KEY AUTOINCREMENT
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
		if err := rows.Scan(&userInput, aiResponse); err != nil {
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

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/query", queryHandler)
	http.HandleFunc("/new-session", newSessionHandler)
	http.HandleFunc("/get-sessions", getSessionsHandler)
	http.HandleFunc("/get-session-history", getSessionHistoryHandler)
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
	// Fetch all session IDs from the sessions table
	rows, err := db.Query("SELECT session_id FROM sessions ORDER BY session_id DESC")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var sessions []int64
	for rows.Next() {
		var sessionID int64
		if err := rows.Scan(&sessionID); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		sessions = append(sessions, sessionID)
	}

	json.NewEncoder(w).Encode(sessions)
}
