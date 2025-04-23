package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/aube/gophermart/internal/ctxkeys"
	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
)

func NewUserWithdrawalsHanlder(
	storeBilling BillingProvider,
	logger *slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userID := ctx.Value(ctxkeys.UserID).(int)

		w.Header().Set("Content-Type", "application/json")

		// Store
		wds, err := storeBilling.Withdrawals(ctx, userID)
		if err != nil {
			logger.ErrorContext(ctx, "UserWithdrawals", "err", err)

			var heherr *httperrors.HTTPError
			if errors.As(err, &heherr) {
				http.Error(w, heherr.Message, heherr.Code)
			} else {
				http.Error(w, "Failed to read user withdrawals", http.StatusInternalServerError)
			}

			return
		}

		if len(wds) == 0 {
			http.Error(w, "No data", http.StatusNoContent)
			return
		}

		// JSON
		result, err := model.WithdrawalsToJSON(wds)

		if err != nil {
			http.Error(w, "Failed to convert user withdrawals", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(result)

		logger.Debug("UserWithdrawals", "User withdrawals", result)
	}
}
