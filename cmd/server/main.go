package main

import (
	"log"

	"github.com/yourname/skribble-backend/internal/config"
	"github.com/yourname/skribble-backend/internal/server"
)

func main() {
	cfg := config.Load()

	srv := server.NewHTTPServer(cfg)

	log.Printf("Starting server on %s\n", cfg.Port)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
