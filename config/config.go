package config

import (
	"glovee-worker/types"
	"os"
	"strconv"
)

func NewConfig() *types.Config {
	return &types.Config{
		Postgres: types.PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     getEnvInt("POSTGRES_PORT", 0),
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DB"),
			SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
		},
		Email: types.EmailConfig{
			ResendAPIKey: os.Getenv("RESEND_API_KEY"),
		},
	}
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
