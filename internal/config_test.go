package app

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		flags          []string
		expectedConfig *envConfig
	}{
		{
			name: "default values",
			expectedConfig: &envConfig{
				ServerAddress:        "localhost:8080",
				LogLevel:             "debug",
				AccrualSystemAddress: "http://localhost:8081",
			},
		},
		{
			name: "env variables override defaults",
			envVars: map[string]string{
				"RUN_ADDRESS":            "env:8080",
				"DATABASE_URI":           "env-dsn",
				"ACCRUAL_SYSTEM_ADDRESS": "env-accrual",
			},
			expectedConfig: &envConfig{
				ServerAddress:        "env:8080",
				LogLevel:             "debug",
				DatabaseDSN:          "env-dsn",
				AccrualSystemAddress: "env-accrual",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Backup and restore original env vars and command line flags
			originalEnv := make(map[string]string)
			for _, e := range os.Environ() {
				originalEnv[e] = os.Getenv(e)
			}
			defer func() {
				for k := range originalEnv {
					os.Unsetenv(k)
				}
				for k, v := range originalEnv {
					os.Setenv(k, v)
				}
			}()

			// Clear env vars
			for k := range originalEnv {
				os.Unsetenv(k)
			}

			// Set test env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Reset command line flags
			flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)

			// Set test flags if provided
			if len(tt.flags) > 0 {
				os.Args = append([]string{"test"}, tt.flags...)
			} else {
				os.Args = []string{"test"}
			}

			config := NewConfig()

			assert.Equal(t, tt.expectedConfig.ServerAddress, config.ServerAddress)
			assert.Equal(t, tt.expectedConfig.LogLevel, config.LogLevel)
			assert.Equal(t, tt.expectedConfig.DatabaseDSN, config.DatabaseDSN)
			assert.Equal(t, tt.expectedConfig.AccrualSystemAddress, config.AccrualSystemAddress)
		})
	}
}
