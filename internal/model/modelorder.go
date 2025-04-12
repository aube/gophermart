package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// User ...
type Order struct {
	ID      int    `json:"id"`
	OrderID int    `json:"order_id"`
	UserID  int    `json:"user_id"`
	Accrual string `json:"accrual"`
	Status  string `json:"status"`
}

// CreateValidate ...
func (o *Order) CreateValidate() error {
	return validation.ValidateStruct(
		o,
		validation.Field(&o.OrderID, validation.Required, is.Int),
		validation.Field(&o.UserID, validation.Required, is.Int),
	)
}

// UpdateValidate ...
func (o *Order) UpdateValidate() error {
	return validation.ValidateStruct(
		o,
		validation.Field(&o.OrderID, validation.Required, is.Int),
	)
}
