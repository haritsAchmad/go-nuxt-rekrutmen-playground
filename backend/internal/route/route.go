package route

import (
	"encoding/json"
	"net/http"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/handler"
)

func RegisterRoutes(lowonganHandler *handler.LowonganHandler, authHandler *handler.AuthHandler) {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/auth/login", authHandler.LoginHandler)
	http.HandleFunc("/api/lowongan", lowonganHandler.LowonganHandler)
	http.HandleFunc("/api/lowongan/detail", lowonganHandler.LowonganDetailHandler)
	http.HandleFunc("/api/lowongan/status", lowonganHandler.LowonganStatusHandler)
	http.HandleFunc("/api/lowongan/bulk-status", lowonganHandler.LowonganBulkStatusHandler)
	http.HandleFunc("/api/lowongan/bulk-delete", lowonganHandler.LowonganBulkDeleteHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Go backend is running",
	})
}
