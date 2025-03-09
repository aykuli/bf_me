package routes

import (
	"bf_me/internal/use_cases"
	"fmt"
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
		sessionId := strings.TrimPrefix(bearerToken, "Token token=")
		fmt.Printf("--------\n\nsession: %s\n\n", sessionId)
		session, _ := uc.Find(sessionId)

		if session != nil {
			next.ServeHTTP(w, r)
		}
		http.Error(w, "Unauthorized", http.StatusForbidden)
	}
}
