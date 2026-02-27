package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type APIConfig struct {
	NATSHost   string
	NATSPort   string
	Port       string
	CORSOrigin string
}

type WorkerConfig struct {
	NATSURL    string
	NATSPort   string
	WorkerName string
}

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment or fallbacks")
	}
}

func requireEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	log.Fatalf("CRITICAL ERROR: Required environment variable '%s' is missing.", key)
	return ""
}

func LoadAPIConfig() *APIConfig {
	Load()
	return &APIConfig{
		NATSHost:   requireEnv("NATS_HOST"),
		NATSPort:   requireEnv("NATS_PORT"),
		Port:       requireEnv("PORT"),
		CORSOrigin: requireEnv("CORS_ORIGIN"),
	}
}

func LoadWorkerConfig() *WorkerConfig {
	Load()
	return &WorkerConfig{
		NATSURL:    requireEnv("NATS_URL"),
		NATSPort:   requireEnv("NATS_PORT"),
		WorkerName: requireEnv("WORKER_NAME"),
	}
}
