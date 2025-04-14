package api

import (
	"errors"
	"io"
	"net/http"

	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
)

func (s *Server) UploadUserOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Body == nil || r.ContentLength == 0 {
		s.logger.ErrorContext(ctx, "UploadUserOrders", "Request body is empty", "")
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.ErrorContext(ctx, "UploadUserOrders", "err", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// OrderID
	order, err := model.ParseOrderID(body)
	if err != nil {
		s.logger.ErrorContext(ctx, "UploadUserOrders", "err", err)
		http.Error(w, "Failed to convert body to uint64", http.StatusBadRequest)
		return
	}

	// Store
	err = s.store.Order.UploadOrders(ctx, &order)
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

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Order uploaded"))

	s.logger.Debug("HandlerCreateUser", "Order uploaded", order.ID)
}

// Загрузка номера заказа
// Хендлер: POST /api/user/orders.
// Хендлер доступен только аутентифицированным пользователям. Номером заказа является последовательность цифр произвольной длины.
// Номер заказа может быть проверен на корректность ввода с помощью алгоритма Луна.
// Формат запроса:

// POST /api/user/orders HTTP/1.1
// Content-Type: text/plain
// ...

// 12345678903

// Возможные коды ответа:

//     200 — номер заказа уже был загружен этим пользователем;
//     202 — новый номер заказа принят в обработку;
//     400 — неверный формат запроса;
//     401 — пользователь не аутентифицирован;
//     409 — номер заказа уже был загружен другим пользователем;
//     422 — неверный формат номера заказа;
//     500 — внутренняя ошибка сервера.
//
//
