package api

import (
	"errors"
	"io"
	"net/http"

	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
)

func (s *Server) UploadUserOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Body == nil || r.ContentLength == 0 {
		s.logger.ErrorContext(ctx, "UploadUserOrders", "Request body is empty", "")
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.ErrorContext(ctx, "UploadUserOrders", "err", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// OrderID
	order, err := model.ParseOrderID(body)
	if err != nil {
		s.logger.ErrorContext(ctx, "UploadUserOrders", "err", err)
		http.Error(w, "Failed to convert body to uint64", http.StatusBadRequest)
		return
	}

	// Store
	err = s.store.Order.UploadOrders(ctx, &order)
	if err != nil {
		s.logger.ErrorContext(ctx, "UploadUserOrders", "err", err)

		var heherr *httperrors.HTTPError
		if errors.As(err, &heherr) {
			http.Error(w, heherr.Message, heherr.Code)
		} else {
			http.Error(w, "Failed to upload order", http.StatusInternalServerError)
		}

		return
	}

	s.store.OrdersQueue.Enqueue(order.ID)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Order uploaded"))

	s.logger.Debug("UploadUserOrders", "Order uploaded", order.ID)
}
