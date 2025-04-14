package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// User ...
type Order struct {
	ID         int    `json:"id"`
	UserID     string `json:"-"`
	Accrual    int    `json:"accrual,omitempty"`
	Status     string `json:"status"`
	UploadedAT string `json:"uploaded_at" db:"created_at"`
}

// CreateValidate ...
func (o *Order) CreateValidate() error {
	return validation.ValidateStruct(
		o,
		validation.Field(&o.UserID, validation.Required, is.Int),
	)
}

// UpdateValidate ...
func (o *Order) UpdateValidate() error {
	return validation.ValidateStruct(
		o,
		validation.Field(&o.ID, validation.Required, is.Int),
	)
}
