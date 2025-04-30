package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Withdraw ...
type Withdraw struct {
	OrderID     int       `json:"order,string"`
	Amount      int64     `json:"-" db:"amount"`
	Sum         float64   `json:"sum" db:"-"`
	ProcessedAt time.Time `json:"processed_at" db:"created_at"`
}

// Validate ...
func (w *Withdraw) Validate() error {
	return validation.ValidateStruct(
		w,
		validation.Field(&w.OrderID, validation.Required),
		validation.Field(&w.Amount, validation.Required),
	)
}

// Balance ...
type Balance struct {
	Withdrawn float64 `json:"withdrawn"`
	Current   float64 `json:"current"`
}
