// internal/room/room.go
// This file defines the Room struct, which represents a game room in the Skribble backend.
// It manages connected clients, handles incoming messages, and maintains the game state for that room.
// The Room struct includes methods for registering and unregistering clients, handling client messages, and broadcasting system messages to all clients in the room.
package room

import (
	"encoding/json"
	"log"

	"github.com/nahom-zewdu/skribble-backend/internal/client"
	"github.com/nahom-zewdu/skribble-backend/internal/game"
	"github.com/nahom-zewdu/skribble-backend/internal/transport"
)

type Room struct {
	ID string

	// Active connected clients in this room
	clients map[string]*client.Client

	// Internal event channels
	register   chan *client.Client
	unregister chan *client.Client
	incoming   chan clientMessage

	// Game state (pure logic container for now)
	game *game.Game
}

// Internal wrapper to preserve sender context
type clientMessage struct {
	client  *client.Client
	message []byte
}

// NewRoom creates a room and starts its event loop.
func NewRoom(id string) *Room {
	r := &Room{
		ID:         id,
		clients:    make(map[string]*client.Client),
		register:   make(chan *client.Client),
		unregister: make(chan *client.Client),
		incoming:   make(chan clientMessage),
		game:       game.NewGame(),
	}

	go r.run()

	return r
}

// Public API (ONLY THESE)
// Register adds a client to the room.
func (r *Room) Register(c *client.Client) {
	r.register <- c
}

// Unregister removes a client from the room.
func (r *Room) Unregister(c *client.Client) {
	r.unregister <- c
}

// HandleClientMessage routes a message from a client into the room.
func (r *Room) HandleClientMessage(c *client.Client, msg []byte) {
	r.incoming <- clientMessage{
		client:  c,
		message: msg,
	}
}

// Room Event Loop
// run is the single authority over room state.
func (r *Room) run() {
	for {
		select {

		case c := <-r.register:
			r.clients[c.ID] = c
			r.onJoin(c)

		case c := <-r.unregister:
			if _, ok := r.clients[c.ID]; ok {
				delete(r.clients, c.ID)
				close(c.Send)
				r.onLeave(c)
			}

		case msg := <-r.incoming:
			r.onMessage(msg.client, msg.message)
		}
	}
}

// Internal Handlers
// onJoin handles logic when a player joins.
func (r *Room) onJoin(c *client.Client) {
	log.Println("Player joined:", c.Name)

	r.broadcastSystem(c.Name + " joined the room")

	// Auto-start when at least 2 players are present
	if len(r.clients) >= 2 && r.game.State == game.Waiting {
		r.game.State = game.Playing
		r.game.CurrentTurn = 1
		r.broadcastSystem("Game started")
	} else if len(r.clients) < 2 {
		r.broadcastSystem("Waiting for players...")
	}
}

// onLeave handles logic when a player disconnects.
func (r *Room) onLeave(c *client.Client) {
	log.Println("Player left:", c.Name)

	r.broadcastSystem(c.Name + " left the room")

	// Pause game if fewer than 2 players remain
	if len(r.clients) < 2 {
		r.game.State = game.Paused
		r.broadcastSystem("Game paused. Waiting for players...")
	}

	// If empty, room manager may clean this later
	if len(r.clients) == 0 {
		log.Println("Room empty:", r.ID)
	}
}

// onMessage handles incoming client messages.
func (r *Room) onMessage(sender *client.Client, raw []byte) {
	var incoming transport.ClientMessage

	if err := json.Unmarshal(raw, &incoming); err != nil {
		log.Println("invalid message format:", err)
		return
	}

	switch incoming.Type {

	case "chat":
		r.handleChat(sender, incoming.Data)

	default:
		log.Println("unknown message type:", incoming.Type)
	}
}

// Internal Utilities
// broadcastSystem sends a system message to all clients.
func (r *Room) broadcastSystem(message string) {
	payload, err := json.Marshal(map[string]interface{}{
		"type": "system",
		"data": message,
	})
	if err != nil {
		log.Println("system message marshal error:", err)
		return
	}

	for _, c := range r.clients {
		c.Send <- payload
	}
}
