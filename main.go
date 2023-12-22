package main

import (
	"GPTChat/db"
	"GPTChat/handlers"
	"log"
	"net/http"
)

func main() {
	dbInstance, err := db.Initialize("./db/conversation.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	handler := handlers.NewHandler(dbInstance)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/log", handler.LogHandler)
	http.HandleFunc("/query", handler.QueryHandler)
	http.HandleFunc("/new-session", handler.NewSessionHandler)
	http.HandleFunc("/get-sessions", handler.GetSessionsHandler)
	http.HandleFunc("/get-session-history", handler.GetSessionHistoryHandler)
	http.HandleFunc("/delete-session", handler.DeleteSessionHandler)
	http.HandleFunc("/rename-session", handler.RenameSessionHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
