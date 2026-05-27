// internal/config/config.go
// This file defines the Config struct, which holds configuration settings for the Skribble backend.
// The Load function reads configuration values from environment variables and returns a Config instance with the appropriate settings.

package config

import "os"

type Config struct {
	Port string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port: ":" + port,
	}
}
