// cmd/server/main.go
// This is the entry point for the Skribble backend server. It loads the configuration, initializes the HTTP server, and starts it.

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nahom-zewdu/skribble-backend/internal/config"
	"github.com/nahom-zewdu/skribble-backend/internal/metrics"
	"github.com/nahom-zewdu/skribble-backend/internal/server"
)

func main() {
	cfg := config.Load()

	metrics.InitRedis(cfg.RedisURL)

	srv := server.NewHTTPServer(cfg)

	log.Printf("Starting server on %s\n", cfg.Port)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	// Graceful shutdown handling
	go func() {
		sig := make(chan os.Signal, 1)

		signal.Notify(
			sig,
			syscall.SIGTERM,
			syscall.SIGINT,
		)

		<-sig

		metrics.Flush()

		os.Exit(0)
	}()
}
