package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresURL         string
	TinybirdHost        string
	TinybirdAppendToken string
	NATSURL             string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Error loading .env file: %v", err)
	}

	return &Config{
		PostgresURL:         getEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/pingpong?sslmode=disable"),
		TinybirdHost:        getEnv("TINYBIRD_HOST", "https://api.tinybird.co"),
		TinybirdAppendToken: getEnv("TINYBIRD_APPEND_TOKEN", ""),
		NATSURL:             getEnv("NATS_URL", "nats://localhost:4222"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
