package api

import (
	"errors"
	"net/http"

	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
)

func (s *Server) UserBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := model.User{}

	// Store
	_, err = s.store.User.Balance(ctx, &user)
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

// Получение текущего баланса пользователя
// Хендлер: GET /api/user/balance.
// Хендлер доступен только авторизованному пользователю. В ответе должны содержаться данные о текущей сумме баллов лояльности, а также сумме использованных за весь период регистрации баллов.
// Формат запроса:

// GET /api/user/balance HTTP/1.1
// Content-Length: 0

// Возможные коды ответа:

//     200 — успешная обработка запроса.
//       Формат ответа:

//   200 OK HTTP/1.1
//   Content-Type: application/json
//   ...

//   {
//       "current": 500.5,
//       "withdrawn": 42
//   }

// 401 — пользователь не авторизован.
// 500 — внутренняя ошибка сервера.
