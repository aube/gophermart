package api

import (
	"errors"
	"net/http"

	"github.com/aube/gophermart/internal/httperrors"
)

func (s *Server) UploadUserOrdersAccrual(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Store
	err := s.store.Order.SetAccrual(ctx, 123123123, 100000)
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

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Order uploaded"))
}
