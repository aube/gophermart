package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// User ...
type Withdraw struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Accrual string `json:"accrual"`
}

// Validate ...
func (w *Withdraw) Validate() error {
	return validation.ValidateStruct(
		w,
		validation.Field(&w.UserID, validation.Required, is.Int),
		validation.Field(&w.Accrual, validation.Required, is.Int),
	)
}
