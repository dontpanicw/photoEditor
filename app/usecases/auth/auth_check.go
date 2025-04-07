package auth

import (
	"fmt"
	"homework-dontpanicw/app/repository"
	"net/http"
	"strconv"
)

func NewAuthMiddleware(repo repository.Session) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionIdString := r.Header.Get("Authorization")
			if sessionIdString == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			sessionId, err := strconv.ParseInt(sessionIdString, 10, 64)
			if err != nil {
				http.Error(w, "Invalid session provided", http.StatusUnauthorized)
				return
			}

			err = repo.GetSession(sessionId)
			if err != nil {
				fmt.Println("sessionId:", sessionId)
				http.Error(w, "Session not found", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
