package errors

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

type CustomError struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	UserMessage string `json:"userMessage"`
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewCustomError(statusCode int, message, userMessage string) *CustomError {
	return &CustomError{
		Code:        statusCode,
		Message:     message,
		UserMessage: userMessage,
	}
}

// HandleError writes an error message to the response and logs the error
func HandleError(w http.ResponseWriter, err error) {
	customErr, ok := err.(*CustomError)
	if !ok {
		customErr = NewCustomError(http.StatusInternalServerError, err.Error(), "An unexpected error has occurred.")
	}

	logError(customErr)
	sendErrorResponse(w, customErr)
}

// logError logs the error to your logging system
func logError(err *CustomError) {
	log.Printf("Error: %s\n", err.Error())
	debug.PrintStack()
}

// sendErrorResponse sends an error response to the client
func sendErrorResponse(w http.ResponseWriter, err *CustomError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(err)
}
