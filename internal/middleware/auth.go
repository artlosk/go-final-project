package middleware

import (
	"go-final-project/internal/auth"
	"go-final-project/internal/helpers"
	"net/http"
	"os"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			helpers.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
			return
		}

		if !auth.ValidateToken(cookie.Value, pass) {
			helpers.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
