package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// User ...
type Order struct {
	ID         int    `json:"id"`
	UserID     int    `json:"-"`
	Accrual    int    `json:"accrual,omitempty"`
	Status     string `json:"status"`
	UploadedAT string `json:"uploaded_at" db:"created_at"`
}

// CreateValidate ...
func (o *Order) UploadOrderValidate() error {
	return validation.ValidateStruct(
		o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.UserID, validation.Required),
	)
}
