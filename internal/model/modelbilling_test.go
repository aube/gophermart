package model

import (
	"testing"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/stretchr/testify/assert"
)

func TestWithdraw_Validate(t *testing.T) {
	tests := []struct {
		name     string
		withdraw Withdraw
		wantErr  bool
	}{
		{
			name: "valid withdraw",
			withdraw: Withdraw{
				OrderID:     744487203422712,
				Amount:      1000,
				Sum:         10.00,
				ProcessedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid order id",
			withdraw: Withdraw{
				OrderID:     0,
				Amount:      1000,
				Sum:         10.00,
				ProcessedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing amount",
			withdraw: Withdraw{
				OrderID:     744487203422712,
				Sum:         10.00,
				ProcessedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.withdraw.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.IsType(t, validation.Errors{}, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBalance_Fields(t *testing.T) {
	balance := Balance{
		Current:   100.50,
		Withdrawn: 50.25,
	}

	// Is float numbers
	assert.Equal(t, 100.50, balance.Current)
	assert.Equal(t, 50.25, balance.Withdrawn)
}
