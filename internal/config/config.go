package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDSN     string
	RedisAddr string
	LimitN    int // Login limit per min
	LimitM    int // Password limit per min
	LimitK    int // IP limit per min
	GRPCPort  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DBDSN:     getEnv("APP_DB_DSN", "postgres://user:password@localhost:5432/antibruteforce?sslmode=disable"),
		RedisAddr: getEnv("APP_REDIS_ADDR", "localhost:6379"),
		GRPCPort:  getEnv("APP_GRPC_PORT", "50051"),
	}

	var err error
	if cfg.LimitN, err = getEnvInt("APP_LIMIT_N", 10); err != nil {
		return nil, err
	}
	if cfg.LimitM, err = getEnvInt("APP_LIMIT_M", 100); err != nil {
		return nil, err
	}
	if cfg.LimitK, err = getEnvInt("APP_LIMIT_K", 1000); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) (int, error) {
	valueStr, ok := os.LookupEnv(key)
	if !ok {
		return fallback, nil
	}
	val, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("env %s must be int: %w", key, err)
	}
	return val, nil
}
