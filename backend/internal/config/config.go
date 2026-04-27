package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ApiConfig struct {
	Port       string
	NATSUrl    string
	CORSOrigin string
}

type WorkerConfig struct {
	NATSUrl    string
	WorkerName string
}

func load() {
	// put variables in os
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment.")
	}
}

func loadFromOS(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		log.Fatalf("Enviroment Variable '%s' missing", key)
		return ""
	}
	return value
}

func LoadApiConfig() *ApiConfig {
	load()
	return &ApiConfig{
		Port:       loadFromOS("PORT"),
		CORSOrigin: loadFromOS("CORS_ORIGIN"),
	}
}

func LoadWorkerConfig() *WorkerConfig {
	load()
	return &WorkerConfig{
		NATSUrl:    loadFromOS("NATS_URL"),
		WorkerName: loadFromOS("WORKER_NAME"),
	}
}
