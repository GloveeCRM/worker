package config

import (
	"os"
	"strconv"
)

type PostgresConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  string
}

type Config struct {
	Postgres PostgresConfig
}

func NewConfig() *Config {
	return &Config{
		Postgres: PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     getEnvInt("POSTGRES_PORT", 0),
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DB"),
			SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
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
