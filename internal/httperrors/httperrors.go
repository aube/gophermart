package httperrors

import (
	"errors"
	"fmt"
	"net/http"
)

type HTTPError struct {
	Code    int
	Err     error
	Message string
}

func (he *HTTPError) Error() string {
	return fmt.Sprintf("%d - %s", he.Code, he.Message)
}

func NewHTTPError(code int, message string) error {
	return &HTTPError{
		Code:    code,
		Message: message,
		Err:     errors.New(message),
	}
}

func NewServerError(err error) error {
	return NewHTTPError(500, fmt.Sprint(err))
}

func NewRecordNotFound() error {
	return NewHTTPError(404, "Record not found")
}

func NewValidationError(validationError error) error {
	return NewHTTPError(http.StatusBadRequest, fmt.Sprint(validationError))
}

func NewConflictError() error {
	return NewHTTPError(http.StatusConflict, "Conflict")
}

func NewLoginFailed() error {
	return NewHTTPError(http.StatusForbidden, "Login failed")
}

func NewAccessDenied() error {
	return NewHTTPError(http.StatusForbidden, "Access denied")
}

func NewAlreadyUploadedError() error {
	return NewHTTPError(http.StatusOK, "Already uploaded")
}
