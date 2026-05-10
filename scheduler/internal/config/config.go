package config

import (
	"os"
	"time"
)

type Config struct {
	PostgresURL string
	NATSURL     string
	BatchSize   int
	TickEvery   time.Duration
}

func Load() Config {
	return Config{
		PostgresURL: envOrDefault("POSTGRES_URL", "postgres://pingpong:pingpong@localhost:5432/pingpong?sslmode=disable"),
		NATSURL:     envOrDefault("NATS_URL", "nats://localhost:4222"),
		BatchSize:   100,
		TickEvery:   time.Second,
	}
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
