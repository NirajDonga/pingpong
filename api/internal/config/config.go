package config

import (
	"errors"
	"os"
)

type Config struct {
	Port          string
	PostgresURL   string
	ClickHouseURL string
	JWTSecret     string
}

func Load() (Config, error) {
	cfg := Config{
		Port:          envOrDefault("PORT", "3001"),
		PostgresURL:   os.Getenv("POSTGRES_URL"),
		ClickHouseURL: os.Getenv("CLICKHOUSE_URL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
	}

	if cfg.PostgresURL == "" {
		return Config{}, errors.New("POSTGRES_URL is required")
	}
	if cfg.ClickHouseURL == "" {
		return Config{}, errors.New("CLICKHOUSE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return Config{}, errors.New("JWT_SECRET is required")
	}

	return cfg, nil
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
