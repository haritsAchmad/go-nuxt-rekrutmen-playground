package handler

import (
	"encoding/json"
	"net/http"
)

type apiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, statusCode int, success bool, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(apiResponse{
		Success: success,
		Message: message,
		Data:    data,
	})
}

func Success(w http.ResponseWriter, message string, data any) {
	writeJSON(w, http.StatusOK, true, message, data)
}

func Created(w http.ResponseWriter, message string, data any) {
	writeJSON(w, http.StatusCreated, true, message, data)
}

func Error(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, false, message, nil)
}
