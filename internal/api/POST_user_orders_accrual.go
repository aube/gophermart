package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/aube/gophermart/internal/httperrors"
)

func NewUploadUserOrdersAccrualHanlder(
	storeOrder OrderProvider,
	logger *slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		// Store
		err := storeOrder.SetAccrual(123123123, 100000)
		if err != nil {
			logger.ErrorContext(ctx, "UploadUserOrders", "err", err)

			var heherr *httperrors.HTTPError
			if errors.As(err, &heherr) {
				http.Error(w, heherr.Message, heherr.Code)
			} else {
				http.Error(w, "Failed to upload order", http.StatusInternalServerError)
			}

			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Order uploaded"))
	}
}
