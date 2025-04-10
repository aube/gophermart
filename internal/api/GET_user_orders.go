package api

import (
	"net/http"
)

func (s *Server) UserOrders(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-token")

	s.logger.Info("ololo")

	if token == "" {
		http.Error(w, "x-token header must be specified", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Получение списка загруженных номеров заказов
// Хендлер: GET /api/user/orders.
// Хендлер доступен только авторизованному пользователю. Номера заказа в выдаче должны быть отсортированы по времени загрузки от самых новых к самым старым. Формат даты — RFC3339.
// Доступные статусы обработки расчётов:

//     NEW — заказ загружен в систему, но не попал в обработку;
//     PROCESSING — вознаграждение за заказ рассчитывается;
//     INVALID — система расчёта вознаграждений отказала в расчёте;
//     PROCESSED — данные по заказу проверены и информация о расчёте успешно получена.

// Формат запроса:

// GET /api/user/orders HTTP/1.1
// Content-Length: 0

// Возможные коды ответа:

//     200 — успешная обработка запроса.
//       Формат ответа:

//   200 OK HTTP/1.1
//   Content-Type: application/json
//   ...

//   [
//       {
//           "number": "9278923470",
//           "status": "PROCESSED",
//           "accrual": 500,
//           "uploaded_at": "2020-12-10T15:15:45+03:00"
//       },
//       {
//           "number": "12345678903",
//           "status": "PROCESSING",
//           "uploaded_at": "2020-12-10T15:12:01+03:00"
//       },
//       {
//           "number": "346436439",
//           "status": "INVALID",
//           "uploaded_at": "2020-12-09T16:09:53+03:00"
//       }
//   ]

// 204 — нет данных для ответа.
// 401 — пользователь не авторизован.
// 500 — внутренняя ошибка сервера.
