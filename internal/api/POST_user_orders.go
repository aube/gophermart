package api

import (
	"io"
	"net/http"

	"github.com/aube/gophermart/internal/model"
)

func (s *Server) UploadUserOrders(w http.ResponseWriter, r *http.Request) {
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
