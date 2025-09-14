package http

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON sends a JSON response
func RespondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// RespondWithError sends an error response
func RespondWithError(w http.ResponseWriter, err error) error {
	w.Header().Set("Content-Type", "application/json")

	// Default to 500 if no specific status code is set
	statusCode := http.StatusInternalServerError

	// You can extend this to handle different error types
	// and set appropriate status codes

	w.WriteHeader(statusCode)

	// Handle nil error
	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
	}

	errorResponse := map[string]string{
		"error": errorMessage,
	}
	return json.NewEncoder(w).Encode(errorResponse)
}

// HandlerFunc is a wrapper for http.HandlerFunc that returns an error
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP implements http.Handler interface
func (hf HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := hf(w, r); err != nil {
		// Log the error or handle it appropriately
		RespondWithError(w, err)
	}
}
