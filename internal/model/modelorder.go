package model

import (
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Order ...
type Order struct {
	ID         int     `json:"number,string"` //lol
	UserID     int     `json:"-"`
	Accrual    float64 `json:"accrual,omitempty"`
	Status     string  `json:"status"`
	UploadedAT string  `json:"uploaded_at" db:"created_at"`
}

// UploadOrderValidate ...
func (o *Order) UploadOrderValidate() error {
	return validation.ValidateStruct(
		o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.UserID, validation.Required),
	)
}

// LuhnCheck ...
func (o *Order) LuhnCheck(orderID int) bool {
	number := strconv.Itoa(orderID)
	sum := 0
	alternate := false

	for i := len(number) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false // contains non-digit character
		}

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = (digit % 10) + 1
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}
