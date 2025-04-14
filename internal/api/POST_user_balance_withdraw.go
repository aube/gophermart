package api

import (
	"errors"
	"io"
	"net/http"

	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
)

func (s *Server) UserBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Body == nil || r.ContentLength == 0 {
		s.logger.ErrorContext(ctx, "UserBalanceWithdraw", "Request body is empty", "")
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.ErrorContext(ctx, "UserBalanceWithdraw", "err", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// JSON
	wd, err := model.ParseWithdraw(body)
	if err != nil {
		s.logger.ErrorContext(ctx, "UserBalanceWithdraw", "err", err)
		return
	}

	user := model.User{}
	// Store
	err = s.store.Billing.BalanceWithdraw(ctx, &wd, &user)
	if err != nil {
		s.logger.ErrorContext(ctx, "UserBalanceWithdraw", "err", err)

		var heherr *httperrors.HTTPError
		if errors.As(err, &heherr) {
			http.Error(w, heherr.Message, heherr.Code)
		} else {
			http.Error(w, "Balance withdraw error", http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ololo, World!"))

}

// Запрос на списание средств
// Хендлер: POST /api/user/balance/withdraw
// Хендлер доступен только авторизованному пользователю. Номер заказа представляет собой гипотетический номер нового заказа пользователя,
// в счёт оплаты которого списываются баллы.
// Примечание: для успешного списания достаточно успешной регистрации запроса,
// никаких внешних систем начисления не предусмотрено и не требуется реализовывать.
// Формат запроса:

// POST /api/user/balance/withdraw HTTP/1.1
// Content-Type: application/json

// {
//     "order": "2377225624",
//     "sum": 751
// }

// Здесь order — номер заказа, а sum — сумма баллов к списанию в счёт оплаты.
// Возможные коды ответа:

//     200 — успешная обработка запроса;
//     401 — пользователь не авторизован;
//     402 — на счету недостаточно средств;
//     422 — неверный номер заказа;
//     500 — внутренняя ошибка сервера.
