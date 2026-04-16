package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPAddr          string
	BackendGraphQLURL string
	RequestTimeout    time.Duration
}

func Load() Config {
	return Config{
		HTTPAddr:          getenv("ACCOUNTS_WEB_HTTP_ADDR", ":8080"),
		BackendGraphQLURL: getenv("ACCOUNTS_WEB_ACCOUNTS_GATEWAY_GRAPHQL_URL", "http://127.0.0.1:8081/query"),
		RequestTimeout:    getDuration("ACCOUNTS_WEB_REQUEST_TIMEOUT", 3*time.Second),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}

	return parsed
}
