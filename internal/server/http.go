// internal/server/http.go
// This file defines the HTTP server for the Skribble backend. It sets up the routes and starts the server.

package server

import (
	"encoding/json"
	"net/http"

	"github.com/nahom-zewdu/skribble-backend/internal/config"
	"github.com/nahom-zewdu/skribble-backend/internal/metrics"
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
	mux.HandleFunc("/metrics", s.metrics)

	server := &http.Server{
		Addr: s.cfg.Port,

		// wrap all routes with middleware
		Handler: s.corsMiddleware(mux),
	}

	return server.ListenAndServe()
}

func (s *HTTPServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		origin := r.Header.Get("Origin")

		allowedOrigins := map[string]bool{
			"http://localhost:5173":         true,
			"https://guessit-nu.vercel.app": true,
		}

		if allowedOrigins[origin] {
			w.Header().Set(
				"Access-Control-Allow-Origin",
				origin,
			)
		}

		w.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, OPTIONS",
		)

		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type",
		)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *HTTPServer) health(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *HTTPServer) metrics(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		metrics.Snapshot(),
	)
}
