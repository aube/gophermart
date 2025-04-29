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
	"github.com/aube/gophermart/internal/model"
	mockApi "github.com/aube/gophermart/mocks/api"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewUserBalanceHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		userID         int
		mockSetup      func(*mockApi.MockUserProvider)
		expectedStatus int
	}{
		{
			name:   "database error",
			userID: 1,
			mockSetup: func(mup *mockApi.MockUserProvider) {
				user := &model.User{ID: 1}
				mup.EXPECT().Balance(gomock.Any(), user).Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserProvider := mockApi.NewMockUserProvider(ctrl)
			tt.mockSetup(mockUserProvider)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := NewUserBalanceHanlder(mockUserProvider, logger)

			req := httptest.NewRequest("GET", "/balance", nil)
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
