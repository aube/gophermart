package api

import (
	"net/http"
)

func (s *Server) UserBalance(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-token")

	s.logger.Info("ololo")

	if token == "" {
		http.Error(w, "x-token header must be specified", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
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
