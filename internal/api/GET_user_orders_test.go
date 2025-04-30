package api

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aube/gophermart/internal/ctxkeys"
	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
	mockApi "github.com/aube/gophermart/mocks/api"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewUserOrdersHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		userID         int
		mockSetup      func(*mockApi.MockOrderProvider)
		expectedStatus int
	}{
		{
			name:   "successful orders retrieval",
			userID: 1,
			mockSetup: func(mop *mockApi.MockOrderProvider) {
				orders := []model.Order{{ID: 123, UserID: 1}}
				mop.EXPECT().Orders(gomock.Any(), 1).Return(orders, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "no orders found",
			userID: 1,
			mockSetup: func(mop *mockApi.MockOrderProvider) {
				mop.EXPECT().Orders(gomock.Any(), 1).Return([]model.Order{}, nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "database error",
			userID: 1,
			mockSetup: func(mop *mockApi.MockOrderProvider) {
				mop.EXPECT().Orders(gomock.Any(), 1).Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "http error",
			userID: 1,
			mockSetup: func(mop *mockApi.MockOrderProvider) {
				err := httperrors.NewHTTPError(http.StatusNotFound, "not found")
				mop.EXPECT().Orders(gomock.Any(), 1).Return(nil, err)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderProvider := mockApi.NewMockOrderProvider(ctrl)
			tt.mockSetup(mockOrderProvider)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := NewUserOrdersHanlder(mockOrderProvider, logger)

			req := httptest.NewRequest("GET", "/orders", nil)
			ctx := context.WithValue(req.Context(), ctxkeys.UserID, tt.userID)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
