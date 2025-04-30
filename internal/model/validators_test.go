package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequiredIf(t *testing.T) {
	tests := []struct {
		name    string
		cond    bool
		value   any
		wantErr bool
	}{
		{
			name:    "condition true with value",
			cond:    true,
			value:   "some value",
			wantErr: false,
		},
		{
			name:    "condition true without value",
			cond:    true,
			value:   "",
			wantErr: true,
		},
		{
			name:    "condition false with value",
			cond:    false,
			value:   "some value",
			wantErr: false,
		},
		{
			name:    "condition false without value",
			cond:    false,
			value:   "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := requiredIf(tt.cond)
			err := rule(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
