package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	RedisURL    string
}

func Load() (*Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return nil, fmt.Errorf("REDIS_URL is required")
	}

	return &Config{
		DatabaseURL: databaseURL,
		RedisURL:    redisURL,
	}, nil
}
