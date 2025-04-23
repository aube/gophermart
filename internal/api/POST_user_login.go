package api

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
)

func NewUserLoginHandler(
	storeUser UserProvider,
	storeActiveUser ActiveUserProvider,
	logger *slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		if r.Body == nil || r.ContentLength == 0 {
			logger.ErrorContext(ctx, "UserLogin", "Request body is empty", "")
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		// Body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.ErrorContext(ctx, "UserLogin", "err", err)
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// JSON
		user, err := model.ParseCredentials(body)
		if err != nil {
			logger.ErrorContext(ctx, "UserLogin", "err", err)
			return
		}

		// Store
		_, err = storeUser.Login(ctx, &user)
		if err != nil {
			logger.ErrorContext(ctx, "UserLogin", "err", err)

			var heherr *httperrors.HTTPError
			if errors.As(err, &heherr) {
				http.Error(w, heherr.Message, heherr.Code)
			} else {
				http.Error(w, "Login failed", http.StatusInternalServerError)
			}

			return
		}

		user.AfterLogin()

		storeActiveUser.Set(ctx, &user)

		setAuthCookie(w, user.RandomHash)

		w.Header().Set("Authorization", bearerString+user.RandomHash)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(user.RandomHash))

		logger.DebugContext(ctx, "UserLogin", "user", user)
	}
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
	c := &http.Cookie{
		Name:     authCookieName,
		Value:    value,
		Expires:  time.Now().Add(24 * time.Hour), // Cookie expires in 24 hours
		Path:     "/",                            // Cookie is accessible across the entire site
		HttpOnly: true,                           // Cookie is not accessible via JavaScript
		Secure:   false,                          // Set to true if using HTTPS
	}

	http.SetCookie(w, c)
}
