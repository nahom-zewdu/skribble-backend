// internal/config/config.go
// This file defines the configuration structure for the Skribble backend server and provides a function to load it.

package config

type Config struct {
	Port string
}

func Load() *Config {
	return &Config{
		Port: ":8080",
	}
}
