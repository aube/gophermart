package api

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/aube/gophermart/internal/ctxkeys"
	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
)

func NewUserBalanceWithdrawHanlder(
	storeUser UserProvider,
	storeBilling BillingProvider,
	logger *slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userID := ctx.Value(ctxkeys.UserID).(int)

		if r.Body == nil || r.ContentLength == 0 {
			logger.ErrorContext(ctx, "UserBalanceWithdraw", "Request body is empty", "")
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		// Body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.ErrorContext(ctx, "UserBalanceWithdraw", "err", err)
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// JSON
		wd, err := model.ParseWithdraw(body)
		if err != nil {
			logger.ErrorContext(ctx, "UserBalanceWithdraw", "err", err)
			return
		}

		user := model.User{
			ID: userID,
		}
		storeUser.Balance(ctx, &user)

		fmt.Println(wd)

		// Store
		err = storeBilling.BalanceWithdraw(ctx, &wd, &user)
		if err != nil {
			logger.ErrorContext(ctx, "UserBalanceWithdraw", "err", err)

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
}
