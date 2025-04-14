package api

import (
	"errors"
	"net/http"

	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
)

func (s *Server) UserOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Store
	orders, err := s.store.Order.Orders(ctx, 313)

	if err != nil {
		s.logger.ErrorContext(ctx, "UploadUserOrders", "err", err)

		var heherr *httperrors.HTTPError
		if errors.As(err, &heherr) {
			http.Error(w, heherr.Message, heherr.Code)
		} else {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		}

		return
	}

	if len(orders) == 0 {
		http.Error(w, "No data", http.StatusNoContent)
		return
	}

	// JSON
	result, err := model.OrdersToJSON(orders)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

	s.logger.Debug("HandlerCreateUser", "Order uploaded", result)
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
