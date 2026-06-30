package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/haritsAchmad/go-nuxt-rekrutmen-playground/backend/internal/auth"
)

type unauthorizedResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Auth(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("Authorization")
			if authorizationHeader == "" {
				writeUnauthorized(w, "Token tidak ditemukan")
				return
			}

			parts := strings.SplitN(authorizationHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				writeUnauthorized(w, "Format token tidak valid")
				return
			}

			claims, err := auth.ValidateToken(parts[1], secretKey)
			if err != nil {
				writeUnauthorized(w, err.Error())
				return
			}

			ctx := auth.WithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	_ = json.NewEncoder(w).Encode(unauthorizedResponse{
		Success: false,
		Message: message,
	})
}
