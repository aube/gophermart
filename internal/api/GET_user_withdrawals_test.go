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

func TestNewUserWithdrawalsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		userID         int
		mockSetup      func(*mockApi.MockBillingProvider)
		expectedStatus int
	}{
		{
			name:   "successful withdrawals retrieval",
			userID: 1,
			mockSetup: func(mbp *mockApi.MockBillingProvider) {
				withdrawals := []model.Withdraw{{OrderID: 123, Sum: 100}}
				mbp.EXPECT().Withdrawals(gomock.Any(), 1).Return(withdrawals, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "no withdrawals found",
			userID: 1,
			mockSetup: func(mbp *mockApi.MockBillingProvider) {
				mbp.EXPECT().Withdrawals(gomock.Any(), 1).Return([]model.Withdraw{}, nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "database error",
			userID: 1,
			mockSetup: func(mbp *mockApi.MockBillingProvider) {
				mbp.EXPECT().Withdrawals(gomock.Any(), 1).Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "http error",
			userID: 1,
			mockSetup: func(mbp *mockApi.MockBillingProvider) {
				err := httperrors.NewHTTPError(http.StatusNotFound, "not found")
				mbp.EXPECT().Withdrawals(gomock.Any(), 1).Return(nil, err)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBillingProvider := mockApi.NewMockBillingProvider(ctrl)
			tt.mockSetup(mockBillingProvider)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := NewUserWithdrawalsHanlder(mockBillingProvider, logger)

			req := httptest.NewRequest("GET", "/withdrawals", nil)
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
