package config

import "os"

type Config struct {
	GRPCAddr string
}

func Load() Config {
	return Config{
		GRPCAddr: getenv("ACCOUNTS_SERVICE_GRPC_ADDR", ":9090"),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
