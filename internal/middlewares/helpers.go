package middlewares

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	// "github.com/aube/url-shortener/internal/app/config"
	// "github.com/aube/url-shortener/internal/logger"
)

const authCookieName = "auth"

func parseJWT(token string, secret string) (*jwt.Token, error) {
	claims := &Claims{}
	return jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

}

func getToken() string {
	var user User
	user.ID = randUserID()

	// Create the JWT claims
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "url-shortener-super-duper-magic-app",
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenSecret := "config.NewConfig().TokenSecret"

	// Sign the token with our secret
	tokenString, err := token.SignedString(tokenSecret)

	// log := logger.Get()

	if err != nil {
		// log.Error("getToken", "token", token)
		// log.Error("getToken", "tokenString", tokenString)
		// log.Error("getToken", "claims", claims)
		return ""
	}
	return tokenString
}

func deleteAuthCookie(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:     authCookieName,
		Value:    "",
		Expires:  time.Unix(0, 0), // Cookie expires in 24 hours
		Path:     "/",             // Cookie is accessible across the entire site
		HttpOnly: true,            // Cookie is not accessible via JavaScript
		Secure:   false,           // Set to true if using HTTPS
	}

	http.SetCookie(w, c)
}

func setAuthCookie(w http.ResponseWriter, value string) {
	// log := logger.Get()
	c := &http.Cookie{
		Name:     authCookieName,
		Value:    value,
		Expires:  time.Now().Add(24 * time.Hour), // Cookie expires in 24 hours
		Path:     "/",                            // Cookie is accessible across the entire site
		HttpOnly: true,                           // Cookie is not accessible via JavaScript
		Secure:   false,                          // Set to true if using HTTPS
	}

	http.SetCookie(w, c)
	// log.Warn("setCookie", "value", value)
}
