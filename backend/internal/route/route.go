package route

import (
	"encoding/json"
	"net/http"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/handler"
)

func RegisterRoutes() {
	http.HandleFunc("/health", withCors(healthHandler))
	http.HandleFunc("/api/lowongan", withCors(handler.LowonganHandler))
	http.HandleFunc("/api/lowongan/detail", withCors(handler.LowonganDetailHandler))
	http.HandleFunc("/api/lowongan/status", withCors(handler.LowonganStatusHandler))
	http.HandleFunc("/api/lowongan/bulk-status", withCors(handler.LowonganBulkStatusHandler))
	http.HandleFunc("/api/lowongan/bulk-delete", withCors(handler.LowonganBulkDeleteHandler))
}

func withCors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Go backend is running",
	})
}