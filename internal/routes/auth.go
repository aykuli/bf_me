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
		sessionId := strings.TrimPrefix(bearerToken, "Token token=")

		session, err := uc.Find(sessionId)
		if session != nil {
			next.ServeHTTP(w, r)
		}

		fmt.Println("session wasnt found", err)
	}
}
