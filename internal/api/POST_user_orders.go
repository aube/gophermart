package api

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/aube/gophermart/internal/ctxkeys"
	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
)

func NewUploadUserOrdersHanlder(
	storeOrder OrderProvider,
	dispatcher AccrualService,
	logger *slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userID := ctx.Value(ctxkeys.UserID).(int)

		if r.Body == nil || r.ContentLength == 0 {
			logger.ErrorContext(ctx, "UploadUserOrders", "Request body is empty", "")
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		// Body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.ErrorContext(ctx, "UploadUserOrders", "err", err)
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// OrderID
		order, err := model.ParseOrderID(body, userID)
		if err != nil {
			logger.ErrorContext(ctx, "UploadUserOrders", "err", err)
			http.Error(w, "Failed to convert body to uint64", http.StatusBadRequest)
			return
		}

		// Store
		err = storeOrder.UploadOrders(ctx, &order)
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

		dispatcher.AddWork(order.ID)

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Order uploaded"))

		logger.Debug("UploadUserOrders", "Order uploaded", order.ID)
	}
}
