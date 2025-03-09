package routes

import (
	"bf_me/internal/use_cases"
	"net/http"
	"strings"
)

func AuthMiddleware(uc *use_cases.SessionsUseCase, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		sessionId := strings.TrimPrefix(bearerToken, "Bearer token=")
		session, err := uc.Find(sessionId)

		if session != nil && err == nil {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "Unauthorized", http.StatusForbidden)
	}
}
