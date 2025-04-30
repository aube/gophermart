package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOrderAccrual(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected OrderAccrual
		wantErr  bool
	}{
		{
			name:  "valid order accrual",
			input: []byte(`{"order": "123", "status": "PROCESSED", "accrual": 10.5}`),
			expected: OrderAccrual{
				ID:      123,
				OrderID: "123",
				Status:  "PROCESSED",
				Sum:     10.5,
				Accrual: 1050,
			},
			wantErr: false,
		},
		{
			name:     "invalid json",
			input:    []byte(`invalid json`),
			expected: OrderAccrual{},
			wantErr:  true,
		},
		{
			name:     "invalid order id",
			input:    []byte(`{"order": "abc", "status": "PROCESSED", "accrual": 10.5}`),
			expected: OrderAccrual{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseOrderAccrual(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestRequest(t *testing.T) {
	tests := []struct {
		name           string
		serverHandler  http.HandlerFunc
		expectedResult OrderAccrual
		expectedError  string
	}{
		{
			name: "successful request",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(OrderAccrual{
					OrderID: "123",
					Status:  "PROCESSED",
					Sum:     10.5,
				})
			},
			expectedResult: OrderAccrual{
				ID:      123,
				OrderID: "123",
				Status:  "PROCESSED",
				Sum:     10.5,
				Accrual: 1050,
			},
		},
		{
			name: "invalid status (204)",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
			expectedError: "invalid",
		},
		{
			name: "too many requests (429)",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTooManyRequests)
			},
			expectedError: "new",
		},
		{
			name: "server error",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedError: "unexpected status code: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.serverHandler)
			defer server.Close()

			result, err := request(server.URL)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}
