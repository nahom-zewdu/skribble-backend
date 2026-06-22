// internal/config/config.go
// This file defines the Config struct and the Load function for loading configuration values from environment variables.
// It includes fields for the server port and Redis connection details.

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	RedisURL string
}

func Load() *Config {
	// Load .env file. It will silently skip if the file doesn't exist (e.g., in production)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environment")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	log.Printf("Config loaded: Port=%s, RedisURL=%s\n", port, redisURL)
	return &Config{
		Port:     ":" + port,
		RedisURL: redisURL,
	}
}
