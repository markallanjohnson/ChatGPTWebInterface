package main

import (
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func main() {
	// Serve static files like index.html
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Handle AJAX requests
	http.HandleFunc("/query", queryHandler)

	// Start the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user query from the request
	userQuery := r.URL.Query().Get("query")

	// Run the Python script with the user query
	cmd := exec.Command("python", "script.py")
	cmd.Stdin = strings.NewReader(userQuery)
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Send the Python script's output back to the client
	w.Write(output)
}
