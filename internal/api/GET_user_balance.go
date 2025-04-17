package api

import (
	"net/http"

	"github.com/aube/gophermart/internal/ctxkeys"
	"github.com/aube/gophermart/internal/model"
)

func (s *Server) UserBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(ctxkeys.UserID).(int)

	user := model.User{
		ID: userID,
	}

	err := s.store.User.Balance(ctx, &user)
	if err != nil {
		http.Error(w, "Failed to get user balance", http.StatusInternalServerError)
		return
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

	s.logger.Debug("UserBalance", "Balance", balance)
}
