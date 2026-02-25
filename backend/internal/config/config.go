package config

import (
	"log"
	"os"

	"github.com/joho/godotenv" // New import
)

// Load reads the .env file and loads it into the environment
func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment or fallbacks")
	}
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
