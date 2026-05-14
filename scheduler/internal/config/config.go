package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresURL string
	NATSURL     string
	BatchSize   int
	TickEvery   time.Duration
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		PostgresURL: mustGetEnv("POSTGRES_URL"),
		NATSURL:     mustGetEnv("NATS_URL"),
		BatchSize:   100,
		TickEvery:   time.Second,
	}
}

func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("CRITICAL STARTUP ERROR: Environment variable '%s' is required but not set.", key)
	}
	return value
}
