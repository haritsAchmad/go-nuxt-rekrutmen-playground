package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				_ = json.NewEncoder(w).Encode(errorResponse{
					Success: false,
					Message: "Terjadi kesalahan pada server",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}
