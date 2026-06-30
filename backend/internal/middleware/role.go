package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/auth"
)

type forbiddenResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	allowedRoleMap := map[string]bool{}

	for _, role := range allowedRoles {
		allowedRoleMap[role] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := auth.ClaimsFromContext(r.Context())
			if !ok {
				writeForbidden(w, "Data user tidak ditemukan")
				return
			}

			if !allowedRoleMap[claims.Role] {
				writeForbidden(w, "Anda tidak punya akses ke aksi ini")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func writeForbidden(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)

	_ = json.NewEncoder(w).Encode(forbiddenResponse{
		Success: false,
		Message: message,
	})
}
