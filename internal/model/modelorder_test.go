package model

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/stretchr/testify/assert"
)

func TestOrder_UploadOrderValidate(t *testing.T) {
	tests := []struct {
		name    string
		order   Order
		wantErr bool
	}{
		{
			name: "valid order",
			order: Order{
				ID:     744487203422712,
				UserID: 1,
			},
			wantErr: false,
		},
		{
			name: "missing order id",
			order: Order{
				UserID: 1,
			},
			wantErr: true,
		},
		{
			name: "missing user id",
			order: Order{
				ID: 744487203422712,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.UploadOrderValidate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.IsType(t, validation.Errors{}, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrder_LuhnCheck(t *testing.T) {
	tests := []struct {
		name    string
		orderID int
		want    bool
	}{
		{
			name:    "valid luhn number",
			orderID: 744487203422712,
			want:    true,
		},
		{
			name:    "invalid luhn number",
			orderID: 12345678902,
			want:    false,
		},
		{
			name:    "zero value",
			orderID: 0,
			want:    true,
		},
		{
			name:    "single digit",
			orderID: 7,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Order{}
			assert.Equal(t, tt.want, o.LuhnCheck(tt.orderID))
		})
	}
}

func TestOrder_Fields(t *testing.T) {
	order := Order{
		ID:         744487203422712,
		UserID:     1,
		Accrual:    100.50,
		Status:     "PROCESSED",
		UploadedAT: "2023-01-01T00:00:00Z",
	}

	assert.Equal(t, 744487203422712, order.ID)
	assert.Equal(t, 1, order.UserID)
	assert.Equal(t, 100.50, order.Accrual)
	assert.Equal(t, "PROCESSED", order.Status)
	assert.Equal(t, "2023-01-01T00:00:00Z", order.UploadedAT)
}
