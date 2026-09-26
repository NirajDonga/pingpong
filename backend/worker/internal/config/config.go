package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	NATSURL    string
	WorkerName string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		NATSURL:    mustGetEnv("NATS_URL"),
		WorkerName: mustGetEnv("WORKER_NAME"),
	}
}

func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("CRITICAL STARTUP ERROR: Environment variable '%s' is required but not set.", key)
	}
	return value
}
