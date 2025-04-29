package api

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
	mockApi "github.com/aube/gophermart/mocks/api"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewUserLoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		body           string
		mockSetup      func(*mockApi.MockUserProvider, *mockApi.MockActiveUserProvider)
		expectedStatus int
	}{
		{
			name:           "empty body",
			body:           "",
			mockSetup:      func(mup *mockApi.MockUserProvider, maup *mockApi.MockActiveUserProvider) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid credentials",
			body: `{"login": "user", "password": "wrong"}`,
			mockSetup: func(mup *mockApi.MockUserProvider, maup *mockApi.MockActiveUserProvider) {
				user := &model.User{Login: "user", Password: "wrong"}
				mup.EXPECT().Login(gomock.Any(), user).Return(nil, httperrors.NewLoginFailed())
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserProvider := mockApi.NewMockUserProvider(ctrl)
			mockActiveUserProvider := mockApi.NewMockActiveUserProvider(ctrl)
			tt.mockSetup(mockUserProvider, mockActiveUserProvider)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := NewUserLoginHandler(mockUserProvider, mockActiveUserProvider, logger)

			req := httptest.NewRequest("POST", "/login", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			handler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
