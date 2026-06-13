# Contributing to Skribble Backend

## Project Overview

Skribble Backend is a Go server that runs a multiplayer drawing and guessing game. Clients connect using WebSockets and join either a public room or a private room.

The backend manages:
- player sessions and room membership
- real-time drawing event broadcast
- game state, turn management, scoring, and timeouts
- hint generation and word selection
- public matchmaking and private room creation/joining

## Architecture

### High-level flow

1. `cmd/server/main.go` loads configuration and starts the HTTP server.
2. `internal/server/http.go` exposes `/ws` for WebSocket upgrades and `/health` for readiness.
3. `internal/server/websocket.go` upgrades the connection, parses query params, and picks the right room.
4. `internal/room/manager.go` creates or returns rooms by ID and supports public/private room lifecycle.
5. `internal/room/room.go` runs the room event loop and coordinates client messages with the game engine.
6. `internal/engine/engine.go` wraps the domain-level game model and exposes time-based ticks.
7. `internal/game/` contains the game domain: players, turns, scoring, timeouts, and word choices.
8. `internal/transport/` defines the message contracts used over WebSockets.

### Core packages

- `internal/config` — configuration loaded from environment variables.
- `internal/server` — HTTP and WebSocket entrypoint logic.
- `internal/client` — WebSocket client read/write pumps.
- `internal/room` — room lifecycle, event dispatch, and broadcast logic.
- `internal/engine` — game loop orchestration.
- `internal/game` — domain logic for gameplay state.
- `internal/transport` — shared message types for server/client communication.
- `internal/pkg/utils` — small utilities such as unique ID generation.

## Key Concepts

### Rooms

A `Room` represents a single game session and keeps:
- connected clients
- the game engine instance
- channels for registration, unregistration, and incoming messages
- a `run()` loop that processes events every 250ms and reacts to client actions

Room types:
- `public` — matchmaking room automatically reused when joinable
- `private` — created explicitly and joined using a room code

A room deletes itself when no clients remain.

### Client lifecycle

- `internal/client/client.go` creates a `Client` with a WebSocket connection.
- `ReadPump()` reads messages and forwards payloads to room handlers.
- `WritePump()` sends server messages back to the client.

### Game engine and domain

The backend separates:
- `engine.Engine` — periodic ticking and command forwarding
- `game.Game` — the actual rules, state machine, and event emission

`game.Game` handles:
- player join / leave
- starting the game once enough players are present
- turn creation and drawer rotation
- word selection, drawing, guessing, scoring
- hints, timeouts, and replay flow

### Messages

#### Incoming message types

- `chat` — guess or chat text from a player
- `select_word` — drawer selects a word
- `draw_start` — begin a stroke
- `draw_move` — continue a stroke
- `draw_end` — finish a stroke
- `clear_canvas` — clear the drawing board

#### Outgoing message types

- `game_snapshot` — initial room/game state
- `system` — system chat messages
- `chat` — player chat/guess messages
- `turn_started` — new turn begins
- `word_selection_started` — drawer chooses from word options
- `drawing_started` — drawing phase begins
- `correct_guess` — guesser guessed correctly
- `turn_ended` — turn complete, next turn scheduled
- `game_ended` — game ends and restart is scheduled
- `selection_timeout` — word selection timed out
- `drawing_timeout` — drawing time expired
- `hint_revealed` — a letter hint was revealed
- `draw_start`, `draw_move`, `draw_end`, `clear_canvas` — drawing broadcast events

## Running the project locally

Prerequisites:
- Go 1.22 or later

Commands:

```bash
go run ./cmd/server
```

Optional port override:

```bash
PORT=9090 go run ./cmd/server
```

Docker:

```bash
docker build -t skribble-backend .
docker run -p 8080:8080 skribble-backend
```

## How to onboard quickly

1. Start with `internal/server/websocket.go` to understand how clients connect.
2. Read `internal/room/manager.go` and `internal/room/room.go` to see room lifecycle and message routing.
3. Follow the game flow in `internal/engine/engine.go` and `internal/game/game.go`.
4. Inspect `internal/game/game_logic.go` for rules, scoring, timeouts, and hints.
5. Review `internal/transport/` to understand the data model used on the socket.

## Common contributor tasks

### Add a new word

Update the word list in `internal/game/word_pool.go`.

### Add a new message type

1. Add the payload type in `internal/transport/`.
2. Update `internal/room/room.go` to parse the incoming message.
3. Add broadcasting logic or state updates as needed.

### Extend game rules

Keep domain rules inside `internal/game/`.
Implement new behavior in `game.Game`, emit events, and consume those events in `internal/room/room.go`.

### Improve rooms or matchmaking

Modify `internal/room/manager.go` for room selection and `internal/room/room.go` for joinability rules.

## Notes for contributors

- Game state is mostly managed inside `internal/game`.
- `internal/room/room.go` is the bridge between socket clients and the domain model.
- Avoid writing raw socket code outside the client/room/server boundary.
- There are no dedicated tests yet, so adding unit tests for game rules is a high-value contribution.

## Useful files

- `cmd/server/main.go` — entry point
- `internal/config/config.go` — app config
- `internal/server/websocket.go` — WebSocket handshake and room selection
- `internal/room/room.go` — message dispatch and broadcasting
- `internal/game/game_logic.go` — game rules, scoring, timeouts
- `internal/game/word_provider.go` — word selection strategy
- `internal/game/word_pool.go` — the static word bank
- `internal/transport/message.go` — socket message contracts

## Recommended first contributions

- Add meaningful unit tests for `internal/game` behavior.
- Improve error handling in WebSocket and room flows.
- Add new categories or harder word choices.
- Add server-side validation for room join logic.
- Document any new message or event contracts in this file.
