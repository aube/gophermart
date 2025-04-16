package api

import (
	"errors"
	"net/http"

	"github.com/aube/gophermart/internal/ctxkeys"
	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
)

func (s *Server) UserOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	// Store
	userID := ctx.Value(ctxkeys.UserID).(int)
	orders, err := s.store.Order.Orders(ctx, userID)

	if err != nil {
		s.logger.ErrorContext(ctx, "UserOrders", "err", err)

		var heherr *httperrors.HTTPError
		if errors.As(err, &heherr) {
			http.Error(w, heherr.Message, heherr.Code)
		} else {
			http.Error(w, "Failed to read user orders", http.StatusInternalServerError)
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
		http.Error(w, "Failed to convert user orders", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)

	s.logger.Debug("UserOrders", "Order uploaded", result)
}
