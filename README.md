# skribble-backend

```txt
skribble-backend/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── transport/
│   │   └── websocket/
│   │       ├── hub.go
│   │       ├── client.go
│   │       └── handlers.go
│   │
│   ├── room/
│   │   ├── manager.go
│   │   ├── room.go
│   │   ├── state.go
│   │   ├── events.go
│   │   └── rotation.go
│   │
│   ├── engine/
│   │   ├── game.go
│   │   ├── scoring.go
│   │   ├── turn.go
│   │   └── word_selector.go
│   │
│   ├── protocol/
│   │   ├── inbound.go
│   │   ├── outbound.go
│   │   └── messages.go
│   │
│   ├── util/
│   │   ├── id.go
│   │   └── slice.go
│   │
│   └── words/
│       └── words.go
│
├── pkg/
│   └── logger/
│       └── logger.go
│
├── go.mod
├── go.sum
└── README.md
```
