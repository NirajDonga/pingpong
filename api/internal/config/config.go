package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	PostgresURL   string
	ClickHouseURL string
	NATSURL       string
	JWTSecret     string
	WebOrigin     string
	CookieSecure  bool
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Port:          mustGetEnv("PORT"),
		PostgresURL:   mustGetEnv("POSTGRES_URL"),
		ClickHouseURL: mustGetEnv("CLICKHOUSE_URL"),
		NATSURL:       mustGetEnv("NATS_URL"),
		JWTSecret:     mustGetEnv("JWT_SECRET"),
		WebOrigin:     mustGetEnv("WEB_ORIGIN"),
		CookieSecure:  mustGetEnvBool("COOKIE_SECURE"),
	}
}

func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("CRITICAL STARTUP ERROR: Environment variable '%s' is required but not set.", key)
	}
	return value
}

func mustGetEnvBool(key string) bool {
	value := mustGetEnv(key)
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("CRITICAL STARTUP ERROR: Environment variable '%s' must be a valid boolean.", key)
	}
	return parsed
}
