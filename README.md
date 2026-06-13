# skribble-backend

**Skribble Backend** is a Go-based WebSocket server that powers a multiplayer drawing and guessing game inspired by Skribbl.io.

## What this repository contains

- `cmd/server/main.go` — application entry point
- `internal/config/config.go` — environment-based runtime configuration
- `internal/server/http.go` — HTTP routing and server startup
- `internal/server/websocket.go` — WebSocket upgrade and session initialization
- `internal/client/client.go` — connected client read/write pumps
- `internal/room/manager.go` — room lifecycle and matchmaking
- `internal/room/room.go` — room event loop, message routing, game integration
- `internal/engine/engine.go` — game loop wrapper and timeout coordination
- `internal/game/` — core gameplay domain model and rules
- `internal/transport/` — WebSocket message definitions and payloads

## Quick start

Run locally:

```bash
go run ./cmd/server
```

Use `PORT` to override the default port:

```bash
PORT=9090 go run ./cmd/server
```

Health check:

```bash
curl http://localhost:8080/health
```

## Documentation

A contributor-friendly guide is available in [`CONTRIBUTING.md`](./CONTRIBUTING.md).
