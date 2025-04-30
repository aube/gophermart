package api

import (
	"bytes"
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

func TestNewUploadUserOrdersHandler(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		userID         int
		mockSetup      func(*mockApi.MockOrderProvider, *mockApi.MockAccrualService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "successful order upload",
			body:   "744487203422712",
			userID: 1,
			mockSetup: func(mop *mockApi.MockOrderProvider, mas *mockApi.MockAccrualService) {
				expectedOrder := model.Order{ID: 744487203422712, UserID: 1}
				mop.EXPECT().UploadOrders(gomock.Any(), &expectedOrder).Return(nil)
				mas.EXPECT().AddWork(744487203422712)
			},
			expectedStatus: http.StatusAccepted,
			expectedBody:   "Order uploaded",
		},
		{
			name:           "empty request body",
			body:           "",
			userID:         1,
			mockSetup:      func(mop *mockApi.MockOrderProvider, mas *mockApi.MockAccrualService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Request body is empty",
		},
		{
			name:           "invalid order number",
			body:           "invalid",
			userID:         1,
			mockSetup:      func(mop *mockApi.MockOrderProvider, mas *mockApi.MockAccrualService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Failed to convert body to uint64",
		},
		{
			name:   "order already uploaded by same user",
			body:   "744487203422712",
			userID: 1,
			mockSetup: func(mop *mockApi.MockOrderProvider, mas *mockApi.MockAccrualService) {
				expectedOrder := model.Order{ID: 744487203422712, UserID: 1}
				err := httperrors.NewHTTPError(http.StatusOK, "order already uploaded by this user")
				mop.EXPECT().UploadOrders(gomock.Any(), &expectedOrder).Return(err)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "order already uploaded by this user",
		},
		{
			name:   "order already uploaded by different user",
			body:   "744487203422712",
			userID: 1,
			mockSetup: func(mop *mockApi.MockOrderProvider, mas *mockApi.MockAccrualService) {
				expectedOrder := model.Order{ID: 744487203422712, UserID: 1}
				err := httperrors.NewHTTPError(http.StatusConflict, "order already uploaded by another user")
				mop.EXPECT().UploadOrders(gomock.Any(), &expectedOrder).Return(err)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   "order already uploaded by another user",
		},
		{
			name:   "internal server error",
			body:   "744487203422712",
			userID: 1,
			mockSetup: func(mop *mockApi.MockOrderProvider, mas *mockApi.MockAccrualService) {
				expectedOrder := model.Order{ID: 744487203422712, UserID: 1}
				mop.EXPECT().UploadOrders(gomock.Any(), &expectedOrder).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to upload order",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOrderProvider := mockApi.NewMockOrderProvider(ctrl)
			mockAccrualService := mockApi.NewMockAccrualService(ctrl)

			tt.mockSetup(mockOrderProvider, mockAccrualService)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := NewUploadUserOrdersHanlder(mockOrderProvider, mockAccrualService, logger)

			req := httptest.NewRequest("POST", "/orders", bytes.NewBufferString(tt.body))
			ctx := context.WithValue(req.Context(), ctxkeys.UserID, tt.userID)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			assert.Contains(t, string(body), tt.expectedBody)
		})
	}
}
