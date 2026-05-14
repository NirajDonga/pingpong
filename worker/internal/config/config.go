package config

import (
	"errors"
	"os"
)

type Config struct {
	NATSURL    string
	WorkerName string
}

func Load() (Config, error) {
	cfg := Config{
		NATSURL:    os.Getenv("NATS_URL"),
		WorkerName: envOrDefault("WORKER_NAME", "worker-local"),
	}

	if cfg.NATSURL == "" {
		return Config{}, errors.New("NATS_URL is required")
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
