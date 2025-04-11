package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// User ...
type Order struct {
	ID            int    `json:"id"`
	OrderID       int    `json:"order_id"`
	UserID        int    `json:"user_id"`
	LoyaltyPoints string `json:"loyalty_points"`
	Status        string `json:"status"`
}

// Validate ...
func (o *Order) Validate() error {
	return validation.ValidateStruct(
		o,
		validation.Field(&o.ID, validation.Required, is.Int),
		validation.Field(&o.UserID, validation.Required, is.Int),
	)
}
