package app

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var configPath = "configs/.env.gophermart"

type envConfig struct {
	ServerAddress         string `envconfig:"RUN_ADDRESS" default:"localhost:8080"`
	LogLevel              string `envconfig:"LOG_LEVEL" default:"debug"`
	DatabaseDSN           string `envconfig:"DATABASE_URI"`
	TokenSecret           string `envconfig:"TOKEN_SECRET"`
	DefaultRequestTimeout string `envconfig:"DEFAULT_REQUEST_TIMEOUT"`
	AccrualSystemAddress  string `envconfig:"ACCRUAL_SYSTEM_ADDRESS"`
}

func NewConfig() *envConfig {
	env := os.Getenv("GO_ENV")

	if env == "" {
		env = ".development"
	}

	godotenv.Load(configPath + env)

	var cfg envConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal("Failed to process env vars: ", err)
	}

	var flagServerAddress string
	var flagDatabaseDSN string
	var flagAccrualSystemAddress string

	flag.StringVar(&flagAccrualSystemAddress, "r", "http://localhost:8081", "accrual system address")
	flag.StringVar(&flagServerAddress, "a", "localhost:8080", "address and port to run service")
	flag.StringVar(&flagDatabaseDSN, "d", "", "Database connection string")
	flag.Parse()

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = flagServerAddress
	}

	if cfg.DatabaseDSN == "" {
		cfg.DatabaseDSN = flagDatabaseDSN
	}

	if cfg.AccrualSystemAddress == "" {
		cfg.AccrualSystemAddress = flagAccrualSystemAddress
	}

	fmt.Println("LogLevel", cfg.LogLevel)
	return &cfg
}
