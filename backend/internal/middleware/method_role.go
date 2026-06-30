package middleware

import "net/http"

func RequireMethodRoles(methodRoles map[string][]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedRoles, ok := methodRoles[r.Method]
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			RequireRoles(allowedRoles...)(next).ServeHTTP(w, r)
		})
	}
}
