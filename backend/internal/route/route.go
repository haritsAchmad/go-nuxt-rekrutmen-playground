package route

import (
	"encoding/json"
	"net/http"

	authhandler "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/handler/auth"
	lowonganhandler "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/handler/lowongan"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/middleware"
)

func RegisterRoutes(lowonganHandler *lowonganhandler.LowonganHandler, authHandler *authhandler.AuthHandler, authSecretKey string) {
	authMiddleware := middleware.Auth(authSecretKey)
	canReadLowongan := middleware.RequireRoles("superadmin", "admin", "viewer")
	canManageLowongan := middleware.RequireRoles("superadmin", "admin")
	canReadOrManageLowongan := middleware.RequireMethodRoles(map[string][]string{
		http.MethodGet:    {"superadmin", "admin", "viewer"},
		http.MethodPost:   {"superadmin", "admin"},
		http.MethodPut:    {"superadmin", "admin"},
		http.MethodDelete: {"superadmin", "admin"},
	})

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/auth/login", authHandler.LoginHandler)
	http.Handle("/api/lowongan", authMiddleware(canReadOrManageLowongan(http.HandlerFunc(lowonganHandler.LowonganHandler))))
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
