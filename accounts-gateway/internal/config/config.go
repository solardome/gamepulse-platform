package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPAddr        string
	AccountGRPCAddr string
	RequestTimeout  time.Duration
}

func Load() Config {
	return Config{
		HTTPAddr:        getenv("ACCOUNTS_GATEWAY_HTTP_ADDR", ":8081"),
		AccountGRPCAddr: getenv("ACCOUNTS_GATEWAY_ACCOUNT_GRPC_ADDR", "127.0.0.1:9090"),
		RequestTimeout:  getDuration("ACCOUNTS_GATEWAY_REQUEST_TIMEOUT", 3*time.Second),
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
