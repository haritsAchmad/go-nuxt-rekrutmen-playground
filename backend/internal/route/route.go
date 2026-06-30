package route

import (
	"encoding/json"
	"net/http"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/handler"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/middleware"
)

func RegisterRoutes(lowonganHandler *handler.LowonganHandler, authHandler *handler.AuthHandler, authSecretKey string) {
	authMiddleware := middleware.Auth(authSecretKey)
	canReadLowongan := middleware.RequireRoles("superadmin", "admin", "viewer")
	canManageLowongan := middleware.RequireRoles("superadmin", "admin")

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/auth/login", authHandler.LoginHandler)
	http.Handle("/api/lowongan", authMiddleware(canReadLowongan(canManageLowongan(http.HandlerFunc(lowonganHandler.LowonganHandler)))))
	http.Handle("/api/lowongan/detail", authMiddleware(canReadLowongan(http.HandlerFunc(lowonganHandler.LowonganDetailHandler))))
	http.Handle("/api/lowongan/status", authMiddleware(canManageLowongan(http.HandlerFunc(lowonganHandler.LowonganStatusHandler))))
	http.Handle("/api/lowongan/bulk-status", authMiddleware(canManageLowongan(http.HandlerFunc(lowonganHandler.LowonganBulkStatusHandler))))
	http.Handle("/api/lowongan/bulk-delete", authMiddleware(canManageLowongan(http.HandlerFunc(lowonganHandler.LowonganBulkDeleteHandler))))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Go backend is running",
	})
}
