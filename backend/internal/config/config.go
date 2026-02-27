package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment or fallbacks")
	}
}

func GetEnv(key string) string {
	if value, _ := os.LookupEnv(key); value != "" {
		return value
	}
	return ""
}

func RequireEnv(key string) string {
	if v := GetEnv(key); v != "" {
		return v
	}
	log.Fatalf("CRITICAL ERROR: Required environment variable '%s' is not set.", key)
	return ""
}
