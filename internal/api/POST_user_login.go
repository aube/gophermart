package api

import (
	"errors"
	"io"
	"net/http"

	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
)

func (s *Server) UserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Body == nil || r.ContentLength == 0 {
		s.logger.ErrorContext(ctx, "UserLogin", "Request body is empty", "")
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.ErrorContext(ctx, "UserLogin", "err", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// JSON
	user, err := model.ParseCredentials(body)
	if err != nil {
		s.logger.ErrorContext(ctx, "UserLogin", "err", err)
		return
	}

	// Store
	_, err = s.store.User.Login(ctx, &user)
	if err != nil {
		s.logger.ErrorContext(ctx, "UserLogin", "err", err)

		var heherr *httperrors.HTTPError
		if errors.As(err, &heherr) {
			http.Error(w, heherr.Message, heherr.Code)
		} else {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ololo, World!"))

}

// Аутентификация пользователя
// Хендлер: POST /api/user/login.
// Аутентификация производится по паре логин/пароль.
// Для передачи аутентификационных данных используйте механизм cookies или HTTP-заголовок Authorization.
// Формат запроса:

// POST /api/user/login HTTP/1.1
// Content-Type: application/json
// ...

// {
//     "login": "<login>",
//     "password": "<password>"
// }
