package api

import (
	"context"
	"net/http"

	"github.com/aube/gophermart/internal/ctxkeys"
)

const authCookieName = "auth"
const bearerString = "Bearer "

func (s *Server) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		authHeader := r.Header.Get("Authorization")

		if authHeader != "" {
			token = authHeader[len(bearerString):]
		}

		if token == "" {
			cookie, err := r.Cookie(authCookieName)
			if err == nil {
				token = cookie.Value
			}
		}

		w.Header().Set("Authorization", bearerString+token)

		user, ok := s.store.ActiveUser.Get(r.Context(), token)

		if !ok {
			http.Error(w, "Auth failed", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxkeys.UserID, user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
