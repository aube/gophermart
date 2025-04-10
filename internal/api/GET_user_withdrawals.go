package api

import (
	"net/http"
)

func (s *Server) UserWithdrawals(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-token")

	s.logger.Info("ololo")

	if token == "" {
		http.Error(w, "x-token header must be specified", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Получение информации о выводе средств
// Хендлер: GET /api/user/withdrawals.
// Хендлер доступен только авторизованному пользователю. Факты выводов в выдаче должны быть отсортированы по времени вывода от самых новых к самым старым. Формат даты — RFC3339.
// Формат запроса:

// GET /api/user/withdrawals HTTP/1.1
// Content-Length: 0

// Возможные коды ответа:

//     200 — успешная обработка запроса.
//       Формат ответа:

//   200 OK HTTP/1.1
//   Content-Type: application/json
//   ...

//   [
//       {
//           "order": "2377225624",
//           "sum": 500,
//           "processed_at": "2020-12-09T16:09:57+03:00"
//       }
//   ]

// 204 — нет ни одного списания.
// 401 — пользователь не авторизован.
// 500 — внутренняя ошибка сервера.
