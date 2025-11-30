package api

import (
	"encoding/json"
	"net/http"
)

// OKResponse writes a JSON payload with a 200 OK status.
func OKResponse(w http.ResponseWriter, data any) {
	jsonResponse(w, http.StatusOK, data)
}

// CreatedResponse writes a JSON payload with a 201 Created status.
func CreatedResponse(w http.ResponseWriter, data any) {
	jsonResponse(w, http.StatusCreated, data)
}

// ErrorResponse writes an error payload with the provided status.
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	if status == 0 {
		status = http.StatusInternalServerError
	}

	response := map[string]string{
		"error": message,
	}

	jsonResponse(w, status, response)
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
