package errors

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

// CustomError represents a structured error for API responses.
type CustomError struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	UserMessage string `json:"userMessage"`
}

// Error implements the error interface for CustomError.
func (e *CustomError) Error() string {
	return e.Message
}

// NewCustomError creates a new instance of CustomError.
func NewCustomError(statusCode int, message, userMessage string) *CustomError {
	return &CustomError{
		Code:        statusCode,
		Message:     message,
		UserMessage: userMessage,
	}
}

// HandleError processes an error, logs it, and sends a corresponding HTTP response.
func HandleError(w http.ResponseWriter, err error) {
	customErr, ok := err.(*CustomError)
	if !ok {
		customErr = NewCustomError(http.StatusInternalServerError, err.Error(), "An unexpected error has occurred.")
	}

	logCustomError(customErr)
	sendErrorResponse(w, customErr)
}

// logCustomError logs the details of the CustomError.
func logCustomError(err *CustomError) {
	log.Printf("Error: %s\n", err.Error())
	debug.PrintStack()
}

// sendErrorResponse sends a structured error response to the client.
func sendErrorResponse(w http.ResponseWriter, err *CustomError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	if encodeErr := json.NewEncoder(w).Encode(err); encodeErr != nil {
		log.Printf("Failed to encode error response: %v", encodeErr)
	}
}
