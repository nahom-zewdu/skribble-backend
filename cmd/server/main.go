// cmd/server/main.go
// This is the entry point for the Skribble backend server. It loads the configuration, initializes the HTTP server, and starts it.

package main

import (
	"log"

	"github.com/nahom-zewdu/skribble-backend/internal/config"
	"github.com/nahom-zewdu/skribble-backend/internal/metrics"
	"github.com/nahom-zewdu/skribble-backend/internal/server"
)

func main() {
	cfg := config.Load()

	metrics.InitRedis(cfg.RedisURL, cfg.RedisToken)

	srv := server.NewHTTPServer(cfg)

	log.Printf("Starting server on %s\n", cfg.Port)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
