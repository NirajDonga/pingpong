package config

import (
	"errors"
	"os"
	"time"
)

type Config struct {
	PostgresURL string
	NATSURL     string
	BatchSize   int
	TickEvery   time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		PostgresURL: os.Getenv("POSTGRES_URL"),
		NATSURL:     os.Getenv("NATS_URL"),
		BatchSize:   100,
		TickEvery:   time.Second,
	}

	if cfg.PostgresURL == "" {
		return Config{}, errors.New("POSTGRES_URL is required")
	}
	if cfg.NATSURL == "" {
		return Config{}, errors.New("NATS_URL is required")
	}

	return cfg, nil
}
