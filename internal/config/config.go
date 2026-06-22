// internal/config/config.go
// This file defines the Config struct and the Load function for loading configuration values from environment variables.
// It includes fields for the server port and Redis connection details.

package config

import "os"

type Config struct {
	Port string

	RedisURL string
}

func Load() *Config {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	return &Config{
		Port: ":" + port,

		RedisURL: os.Getenv(
			"Redis_URL",
		),
	}
}
