package route

import (
	"encoding/json"
	"net/http"

	authhandler "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/handler/auth"
	lowonganhandler "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/handler/lowongan"
	pelamarhandler "github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/handler/pelamar"
	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/middleware"
)

func RegisterRoutes(lowonganHandler *lowonganhandler.LowonganHandler, pelamarHandler *pelamarhandler.PelamarHandler, authHandler *authhandler.AuthHandler, authSecretKey string) {
	authMiddleware := middleware.Auth(authSecretKey)
	canReadLowongan := middleware.RequireRoles("superadmin", "admin", "viewer")
	canManageLowongan := middleware.RequireRoles("superadmin", "admin")
	canReadOrManageLowongan := middleware.RequireMethodRoles(map[string][]string{
		http.MethodGet:    {"superadmin", "admin", "viewer"},
		http.MethodPost:   {"superadmin", "admin"},
		http.MethodPut:    {"superadmin", "admin"},
		http.MethodDelete: {"superadmin", "admin"},
	})
	canReadPelamar := middleware.RequireRoles("superadmin", "admin", "viewer")
	canManagePelamar := middleware.RequireRoles("superadmin", "admin")
	canReadOrManagePelamar := middleware.RequireMethodRoles(map[string][]string{
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
	http.Handle("/api/pelamar", authMiddleware(canReadOrManagePelamar(http.HandlerFunc(pelamarHandler.PelamarHandler))))
	http.Handle("/api/pelamar/detail", authMiddleware(canReadPelamar(http.HandlerFunc(pelamarHandler.PelamarDetailHandler))))
	http.Handle("/api/pelamar/status", authMiddleware(canManagePelamar(http.HandlerFunc(pelamarHandler.PelamarStatusHandler))))
	http.Handle("/api/pelamar/bulk-status", authMiddleware(canManagePelamar(http.HandlerFunc(pelamarHandler.PelamarBulkStatusHandler))))
	http.Handle("/api/pelamar/bulk-delete", authMiddleware(canManagePelamar(http.HandlerFunc(pelamarHandler.PelamarBulkDeleteHandler))))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Go backend is running",
	})
}
