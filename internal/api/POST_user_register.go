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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered"))

	s.logger.Debug("HandlerCreateUser", "User registered", user.ID)
}

// Регистрация пользователя
// Хендлер: POST /api/user/register.
// Регистрация производится по паре логин/пароль. Каждый логин должен быть уникальным.
// После успешной регистрации должна происходить автоматическая аутентификация пользователя.
// Для передачи аутентификационных данных используйте механизм cookies или HTTP-заголовок Authorization.
// Формат запроса:

// POST /api/user/register HTTP/1.1
// Content-Type: application/json
// ...

// {
//     "login": "<login>",
//     "password": "<password>"
// }

// Возможные коды ответа:

//     200 — пользователь успешно зарегистрирован и аутентифицирован;
//     400 — неверный формат запроса;
//     409 — логин уже занят;
//     500 — внутренняя ошибка сервера.
