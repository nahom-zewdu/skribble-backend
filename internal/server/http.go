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

// Start initializes the HTTP server and listens for incoming requests on the configured port.
func (s *HTTPServer) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", s.handleWebSocket)
	mux.HandleFunc("/health", s.health)

	server := &http.Server{
		Addr:    s.cfg.Port,
		Handler: mux,
	}

	return server.ListenAndServe()
}

func (s *HTTPServer) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
