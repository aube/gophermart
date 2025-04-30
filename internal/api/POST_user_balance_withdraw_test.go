package api

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aube/gophermart/internal/ctxkeys"
	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
	mockApi "github.com/aube/gophermart/mocks/api"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewUserBalanceWithdrawHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		userID         int
		body           string
		mockSetup      func(*mockApi.MockUserProvider, *mockApi.MockBillingProvider)
		expectedStatus int
	}{
		{
			name:   "successful withdraw",
			userID: 1,
			body:   `{"order": "12345678903", "sum": 100}`,
			mockSetup: func(mup *mockApi.MockUserProvider, mbp *mockApi.MockBillingProvider) {
				user := &model.User{ID: 1}
				mup.EXPECT().Balance(gomock.Any(), user).Return(nil)
				mbp.EXPECT().BalanceWithdraw(gomock.Any(), gomock.Any(), user).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty body",
			userID:         1,
			body:           "",
			mockSetup:      func(mup *mockApi.MockUserProvider, mbp *mockApi.MockBillingProvider) {},
			expectedStatus: http.StatusBadRequest,
		},
		// {
		// 	name:   "invalid body",
		// 	userID: 1,
		// 	body:   `invalid json`,
		// 	mockSetup: func(mup *mockApi.MockUserProvider, mbp *mockApi.MockBillingProvider) {
		// 		user := &model.User{ID: 1}
		// 		mup.EXPECT().Balance(gomock.Any(), user).Return(nil)
		// 	},
		// 	expectedStatus: http.StatusBadRequest,
		// },
		{
			name:   "not enough money",
			userID: 1,
			body:   `{"order": "12345678903", "sum": 100}`,
			mockSetup: func(mup *mockApi.MockUserProvider, mbp *mockApi.MockBillingProvider) {
				user := &model.User{ID: 1}
				mup.EXPECT().Balance(gomock.Any(), user).Return(nil)
				err := httperrors.NewNotEnoughMoneyError()
				mbp.EXPECT().BalanceWithdraw(gomock.Any(), gomock.Any(), user).Return(err)
			},
			expectedStatus: 402,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserProvider := mockApi.NewMockUserProvider(ctrl)
			mockBillingProvider := mockApi.NewMockBillingProvider(ctrl)
			tt.mockSetup(mockUserProvider, mockBillingProvider)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := NewUserBalanceWithdrawHanlder(mockUserProvider, mockBillingProvider, logger)

			req := httptest.NewRequest("POST", "/withdraw", strings.NewReader(tt.body))
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
