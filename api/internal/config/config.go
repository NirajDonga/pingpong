package config

import (
	"errors"
	"os"
)

type Config struct {
	Port        string
	PostgresDSN string
	JWTSecret   string
}

func Load() (Config, error) {
	cfg := Config{
		Port:        envOrDefault("PORT", "3001"),
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	if cfg.PostgresDSN == "" {
		return Config{}, errors.New("POSTGRES_DSN is required")
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
