package api

import (
	"net/http"

	"github.com/aube/gophermart/internal/model"
)

func (s *Server) UserBalance(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	user := model.User{
		Balance:   111111,
		Withdrawn: 111111,
	}

	balance := model.Balance{
		Current:   float64(user.Balance) / 100,
		Withdrawn: float64(user.Withdrawn) / 100,
	}

	json, err := model.BalanceToJSON(balance)

	if err != nil {
		http.Error(w, "Failed to convert user balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
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
