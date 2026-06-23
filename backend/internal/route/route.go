package route

import (
	"encoding/json"
	"net/http"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/handler"
)

func RegisterRoutes() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/lowongan", handler.LowonganHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Go backend is running",
	})
}