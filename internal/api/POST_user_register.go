package api

import (
	"io"
	"net/http"

	"github.com/aube/gophermart/internal/model"
)

func (s *Server) UserRegister(w http.ResponseWriter, r *http.Request) {
	httpStatus := http.StatusCreated

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
		return
	}

	// Store
	err = s.store.User.Register(ctx, &user)
	if err != nil {
		s.logger.ErrorContext(ctx, "HandlerCreateUser", "err", err)
		httpStatus = http.StatusConflict
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write([]byte("Ololo, World!"))

	s.logger.Debug("HandlerCreateUser", "httpStatus", err)
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
