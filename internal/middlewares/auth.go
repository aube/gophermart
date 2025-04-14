package middlewares

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"

	// "github.com/aube/url-shortener/internal/app/config"
	"github.com/aube/gophermart/internal/ctxkeys"
	// "github.com/aube/url-shortener/internal/logger"
)

const bearerString = "Bearer " // The token should be in the format "Bearer <token>"

type Claims struct {
	UserID string `json:"id"`
	jwt.RegisteredClaims
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ID       string `json:"id"`
}

func randUserID() string {
	min := 11111
	max := 99999
	rndInt := rand.Intn(max-min) + min
	return strconv.Itoa(rndInt)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log := logger.WithContext(r.Context())

		var tokenString string

		authHeader := r.Header.Get("Authorization")

		if authHeader != "" {
			tokenString = authHeader[len(bearerString):]
		}

		if tokenString == "" {
			cookie, err := r.Cookie(authCookieName)
			if err == nil {
				tokenString = cookie.Value
			}
		}

		w.Header().Set("Authorization", bearerString+tokenString)

		token, err := parseJWT(tokenString, "config.NewConfig().TokenSecret")

		if err != nil {
			// log.Error("AuthMiddleware", "err", err)
			// log.Warn("AuthMiddleware", "token", token)

			if err == jwt.ErrSignatureInvalid {
				tokenErrorMsg = "Invalid token signature"
			} else {
				tokenErrorMsg = "Invalid token"
			}
		}

		if !token.Valid {
			tokenErrorMsg = "Invalid token"
		}

		if tokenErrorMsg != "" {
			deleteAuthCookie(w)
			http.Error(w, tokenErrorMsg, http.StatusUnauthorized)
			return
		}

		setAuthCookie(w, bearerString+tokenString)
		UserID := claims.UserID

		ctx := context.WithValue(r.Context(), ctxkeys.UserIDKey, UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
