package api

import (
	"errors"
	"io"
	"net/http"

	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
)

func (s *Server) UserRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Body == nil || r.ContentLength == 0 {
		s.logger.ErrorContext(ctx, "HandlerCreateUser", "Request body is empty", "")
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.ErrorContext(ctx, "HandlerCreateUser", "err", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// JSON
	user, err := model.ParseCredentials(body)
	if err != nil {
		s.logger.ErrorContext(ctx, "HandlerCreateUser", "err", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Store
	err = s.store.User.Register(ctx, &user)
	if err != nil {
		s.logger.ErrorContext(ctx, "HandlerCreateUser", "err", err)

		var heherr *httperrors.HTTPError
		if errors.As(err, &heherr) {
			http.Error(w, heherr.Message, heherr.Code)
		} else {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		}

		return
	}

	user.AfterLogin()

	s.store.ActiveUser.Set(ctx, &user)

	setAuthCookie(w, user.RandomHash)

	w.Header().Set("Authorization", bearerString+user.RandomHash)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered"))

	s.logger.Debug("HandlerCreateUser", "User registered", user.ID)
}
