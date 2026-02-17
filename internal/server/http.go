// internal/server/http.go
// This file defines the HTTP server for the Skribble backend. It sets up the routes and starts the server.

package server

import (
	"net/http"

	"github.com/nahom-zewdu/skribble-backend/internal/config"
)

type HTTPServer struct {
	cfg *config.Config
}

func NewHTTPServer(cfg *config.Config) *HTTPServer {
	return &HTTPServer{
		cfg: cfg,
	}
}

func (s *HTTPServer) Start() error {
	http.HandleFunc("/ws", s.handleWebSocket)
	return http.ListenAndServe(s.cfg.Port, nil)
}
