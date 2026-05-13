package config

import "os"

type Config struct {
	NATSURL    string
	WorkerName string
}

func Load() Config {
	return Config{
		NATSURL:    envOrDefault("NATS_URL", "nats://localhost:4222"),
		WorkerName: envOrDefault("WORKER_NAME", "worker-local"),
	}
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
