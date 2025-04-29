package httperrors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPError(t *testing.T) {
	tests := []struct {
		name        string
		code        int
		message     string
		expectedErr string
	}{
		{
			name:        "bad request",
			code:        http.StatusBadRequest,
			message:     "Invalid input",
			expectedErr: "400 - Invalid input",
		},
		{
			name:        "not found",
			code:        http.StatusNotFound,
			message:     "Not found",
			expectedErr: "404 - Not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewHTTPError(tt.code, tt.message)
			assert.Equal(t, tt.expectedErr, err.Error())
		})
	}

	t.Run("assertion type check it's an HTTPError", func(t *testing.T) {
		err := NewHTTPError(tests[0].code, tests[0].message)
		assert.Equal(t, tests[0].expectedErr, err.Error())

		httpErr, ok := err.(*HTTPError)
		assert.True(t, ok)
		assert.Equal(t, tests[0].code, httpErr.Code)
		assert.Equal(t, tests[0].message, httpErr.Message)
	})
}

func TestNewServerError(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedErr string
	}{
		{
			name:        "standard error",
			inputErr:    errors.New("database connection failed"),
			expectedErr: "500 - database connection failed",
		},
		{
			name:        "nil error",
			inputErr:    nil,
			expectedErr: "500 - <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewServerError(tt.inputErr)
			assert.Equal(t, tt.expectedErr, err.Error())
		})
	}
}

func TestSpecializedErrorFunctions(t *testing.T) {
	tests := []struct {
		name         string
		fn           func() error
		expectedErr  string
		expectedCode int
	}{
		{
			name:         "RecordNotFound",
			fn:           NewRecordNotFound,
			expectedErr:  "404 - Record not found",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "ValidationError",
			fn:           func() error { return NewValidationError(errors.New("validation failed")) },
			expectedErr:  "400 - validation failed",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "ConflictError",
			fn:           NewConflictError,
			expectedErr:  "409 - Conflict",
			expectedCode: http.StatusConflict,
		},
		{
			name:         "LoginFailed",
			fn:           NewLoginFailed,
			expectedErr:  "403 - Login failed",
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "AccessDenied",
			fn:           NewAccessDenied,
			expectedErr:  "403 - Access denied",
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "AlreadyUploadedByMeError",
			fn:           NewAlreadyUploadedByMeError,
			expectedErr:  "200 - Already uploaded",
			expectedCode: http.StatusOK,
		},
		{
			name:         "AlreadyUploadedAnotherError",
			fn:           NewAlreadyUploadedAnotherError,
			expectedErr:  "409 - Already uploaded",
			expectedCode: http.StatusConflict,
		},
		{
			name:         "NotEnoughMoneyError",
			fn:           NewNotEnoughMoneyError,
			expectedErr:  "402 - I need more gold!",
			expectedCode: 402,
		},
		{
			name:         "OrderNumberError",
			fn:           NewOrderNumberError,
			expectedErr:  "422 - Order number incorrect",
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			assert.Equal(t, tt.expectedErr, err.Error())

			httpErr, ok := err.(*HTTPError)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, httpErr.Code)
		})
	}
}
