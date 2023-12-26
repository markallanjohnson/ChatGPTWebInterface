package main

import (
	"GPTChat/db"
	"GPTChat/handlers"
	"GPTChat/services"
	"log"
	"net/http"
)

func main() {
	// Initialize the database connection
	dbInstance, err := db.Initialize("./db/conversation.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create repositories
	sessionRepo := db.NewSessionRepository(dbInstance.SQLDB())
	historyRepo := db.NewHistoryRepository(dbInstance.SQLDB())

	// Create services with the repos
	chatService := services.NewChatService(historyRepo)
	sessionService := services.NewSessionService(sessionRepo)

	// Create a new handler with dependencies
	handler := handlers.NewHandler(chatService, sessionService)

	// Set up HTTP routes
	setupRoutes(handler)

	// Start the server
	startServer()
}

// setupRoutes configures the HTTP routes for the server.
func setupRoutes(handler *handlers.Handler) {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/log", handler.LogHandler)
	http.HandleFunc("/query", handler.QueryHandler)
	http.HandleFunc("/new-session", handler.NewSessionHandler)
	http.HandleFunc("/get-sessions", handler.GetSessionsHandler)
	http.HandleFunc("/get-session-history", handler.GetSessionHistoryHandler)
	http.HandleFunc("/delete-session", handler.DeleteSessionHandler)
	http.HandleFunc("/rename-session", handler.RenameSessionHandler)
}

// startServer starts the HTTP server.
func startServer() {
	const address = ":8080"
	log.Printf("Starting server on %s", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
