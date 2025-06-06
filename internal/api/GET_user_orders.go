package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/aube/gophermart/internal/ctxkeys"
	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
)

func NewUserOrdersHanlder(
	storeOrder OrderProvider,
	logger *slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userID := ctx.Value(ctxkeys.UserID).(int)

		w.Header().Set("Content-Type", "application/json")

		// Store
		orders, err := storeOrder.Orders(ctx, userID)

		if err != nil {
			logger.ErrorContext(ctx, "UserOrders", "err", err)

			var heherr *httperrors.HTTPError
			if errors.As(err, &heherr) {
				http.Error(w, heherr.Message, heherr.Code)
			} else {
				http.Error(w, "Failed to read user orders", http.StatusInternalServerError)
			}

			return
		}

		if len(orders) == 0 {
			w.WriteHeader(http.StatusNoContent)
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

		logger.Debug("UserOrders", "User orders", orders)
	}
}
