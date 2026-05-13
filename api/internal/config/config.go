package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	Port          string
	PostgresURL   string
	ClickHouseURL string
	NATSURL       string
	JWTSecret     string
	WebOrigin     string
	CookieSecure  bool
}

func Load() (Config, error) {
	cfg := Config{
		Port:          envOrDefault("PORT", "3001"),
		PostgresURL:   os.Getenv("POSTGRES_URL"),
		ClickHouseURL: os.Getenv("CLICKHOUSE_URL"),
		NATSURL:       os.Getenv("NATS_URL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		WebOrigin:     envOrDefault("WEB_ORIGIN", "http://localhost:3000"),
		CookieSecure:  envBoolOrDefault("COOKIE_SECURE", false),
	}

	if cfg.PostgresURL == "" {
		return Config{}, errors.New("POSTGRES_URL is required")
	}
	if cfg.ClickHouseURL == "" {
		return Config{}, errors.New("CLICKHOUSE_URL is required")
	}
	if cfg.NATSURL == "" {
		return Config{}, errors.New("NATS_URL is required")
	}
	if cfg.JWTSecret == "" {
		return Config{}, errors.New("JWT_SECRET is required")
	}

	return cfg, nil
}

func envBoolOrDefault(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
