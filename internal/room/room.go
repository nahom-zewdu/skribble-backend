// internal/room/room.go
// This file defines the Room struct, which represents a game room in the Skribble backend.
// It includes fields for the room ID, a map of connected clients, channels for registering and unregistering clients, and a channel for broadcasting messages to all clients in the room. The Room struct also has a method to run the main loop for handling client connections and messages.
package room

import (
	"log"

	"github.com/nahom-zewdu/skribble-backend/internal/client"
	"github.com/nahom-zewdu/skribble-backend/internal/game"
)

type Room struct {
	ID string

	clients map[string]*client.Client

	register   chan *client.Client
	unregister chan *client.Client
	broadcast  chan []byte

	game *game.Game
}

func NewRoom(id string) *Room {
	r := &Room{
		ID:         id,
		clients:    make(map[string]*client.Client),
		register:   make(chan *client.Client),
		unregister: make(chan *client.Client),
		broadcast:  make(chan []byte),
		game:       game.NewGame(),
	}

	go r.run()

	return r
}

func (r *Room) run() {
	// Main loop to handle client registration, unregistration, and message broadcasting
	for {
		select {

		case c := <-r.register:
			r.clients[c.ID] = c
			r.handleJoin(c)

		case c := <-r.unregister:
			if _, ok := r.clients[c.ID]; ok {
				delete(r.clients, c.ID)
				close(c.Send)
				r.handleLeave(c)
			}

		case message := <-r.broadcast:
			for _, c := range r.clients {
				select {
				case c.Send <- message:
				default:
					close(c.Send)
					delete(r.clients, c.ID)
				}
			}
		}
	}
}

func (r *Room) handleJoin(c *client.Client) {
	// Handle a new client joining the room
	log.Println("Player joined:", c.Name)

	if len(r.clients) >= 2 && r.game.State == game.Waiting {
		r.game.State = game.Playing
		r.game.CurrentTurn = 1
		r.broadcastSystem("Game started")
	} else {
		r.broadcastSystem("Waiting for players...")
	}
}
