// internal/room/room.go
// This file defines the Room struct, which manages a single game room in the Skribble backend.
// It handles client registration, message processing, and game state management for the room.
// The Room struct interacts with the Game logic to manage turns, scoring, and player actions.
package room

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/nahom-zewdu/skribble-backend/internal/client"
	"github.com/nahom-zewdu/skribble-backend/internal/engine"
	"github.com/nahom-zewdu/skribble-backend/internal/game"
	"github.com/nahom-zewdu/skribble-backend/internal/transport"
)

type Room struct {
	ID string

	clients map[string]*client.Client

	register   chan *client.Client
	unregister chan *client.Client
	incoming   chan clientMessage

	game   *game.Game
	engine *engine.Engine
}

type clientMessage struct {
	client  *client.Client
	message []byte
}

func NewRoom(id string) *Room {
	g := game.NewGame()
	e := engine.New(g)

	r := &Room{
		ID:         id,
		clients:    make(map[string]*client.Client),
		register:   make(chan *client.Client),
		unregister: make(chan *client.Client),
		incoming:   make(chan clientMessage),
		game:       g,
		engine:     e,
	}

	go r.run()
	return r
}

func (r *Room) Register(c *client.Client) {
	r.register <- c
}

func (r *Room) Unregister(c *client.Client) {
	r.unregister <- c
}

func (r *Room) HandleClientMessage(c *client.Client, msg []byte) {
	r.incoming <- clientMessage{client: c, message: msg}
}

// run is the main loop for the room, handling client registration, unregistration, and incoming messages.
func (r *Room) run() {
	for {
		select {

		case c := <-r.register:
			r.clients[c.ID] = c
			r.game.AddPlayer(c.ID, c.Name)
			r.onJoin(c)

		case c := <-r.unregister:
			if _, ok := r.clients[c.ID]; ok {
				delete(r.clients, c.ID)
				close(c.Send)
				r.game.RemovePlayer(c.ID)
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
	r.broadcastSystem(c.Name + " joined")

	if len(r.clients) >= 2 && r.game.State == game.Waiting {
		if _, err := r.game.Start(); err == nil {
			r.startTurnBroadcast()
		}
	}
}

// onLeave handles logic when a player disconnects.
func (r *Room) onLeave(c *client.Client) {
	r.broadcastSystem(c.Name + " left")

	if r.game.CurrentTurn != nil && c.ID == r.game.CurrentTurn.DrawerID {
		r.endTurn()
	}

	if len(r.clients) < 2 {
		r.game.State = game.Paused
		r.broadcastSystem("Game paused")
	}
}

// onMessage processes incoming messages from clients and routes them to the appropriate handlers.
func (r *Room) onMessage(sender *client.Client, raw []byte) {
	var incoming transport.ClientMessage

	if err := json.Unmarshal(raw, &incoming); err != nil {
		return
	}

	switch incoming.Type {

	case "chat":
		r.handleChat(sender, incoming.Data)
	}
}

// handleChat processes chat messages and checks for correct guesses.
func (r *Room) handleChat(sender *client.Client, raw json.RawMessage) {
	var chat transport.ChatMessage
	if err := json.Unmarshal(raw, &chat); err != nil {
		return
	}

	// Guess logic
	if r.game.CurrentTurn != nil &&
		strings.EqualFold(chat.Text, r.game.CurrentTurn.Word) {

		_, err := r.game.Guess(sender.ID, chat.Text)
		if err != nil {
			log.Println("Guess error:", err)
			return
		}
		if r.game.CurrentTurn.Guessed[sender.ID] == true {
			r.broadcastCorrectGuess(sender.Name)

			if r.allGuessed() {
				r.endTurn()
			}
			return
		}
	}

	// Normal chat broadcast
	r.broadcastChat(sender.Name, chat.Text)
}

// Helper methods for broadcasting messages and managing game state
func (r *Room) allGuessed() bool {
	turn := r.game.CurrentTurn
	if turn == nil {
		return false
	}

	for _, p := range r.game.Players {
		if p.ID == turn.DrawerID {
			continue
		}
		if !turn.Guessed[p.ID] {
			return false
		}
	}
	return true
}

// endTurn ends the current turn and starts the next one, broadcasting updates to clients.
func (r *Room) endTurn() {
	if _, err := r.game.EndTurn(); err != nil {
		r.broadcastSystem("Game ended")
		return
	}

	r.broadcastSystem("Turn ended")
	r.startTurnBroadcast()
}

// startTurnBroadcast sends a message to all clients indicating the start of a new turn and who the drawer is.
func (r *Room) startTurnBroadcast() {
	turn := r.game.CurrentTurn
	if turn == nil {
		return
	}

	r.broadcastSystem("Turn " + string(rune(turn.Number+'0')) + " started")

	for _, c := range r.clients {
		if c.ID == turn.DrawerID {
			c.Send <- r.mustMarshal(transport.Message{
				Type: "your_turn",
				Data: map[string]interface{}{
					"turn": turn.Number,
				},
			})
		} else {
			c.Send <- r.mustMarshal(transport.Message{
				Type: "turn_started",
				Data: map[string]interface{}{
					"turn": turn.Number,
				},
			})
		}
	}
}

// broadcastChat sends a chat message from a player to all clients in the room.
func (r *Room) broadcastChat(sender, text string) {
	msg := transport.Message{
		Type: "chat",
		Data: map[string]interface{}{
			"sender": sender,
			"text":   text,
		},
	}

	payload := r.mustMarshal(msg)

	for _, c := range r.clients {
		c.Send <- payload
	}
}

// broadcastCorrectGuess sends a message to all clients indicating that a player has made a correct guess.
func (r *Room) broadcastCorrectGuess(name string) {
	msg := transport.Message{
		Type: "correct_guess",
		Data: map[string]interface{}{
			"player": name,
		},
	}

	payload := r.mustMarshal(msg)

	for _, c := range r.clients {
		c.Send <- payload
	}
}

// broadcastSystem sends a system message to all clients in the room, typically used for notifications like players joining/leaving or game state changes.
func (r *Room) broadcastSystem(text string) {
	msg := transport.Message{
		Type: "system",
		Data: transport.SystemMessage{Text: text},
	}

	payload := r.mustMarshal(msg)

	for _, c := range r.clients {
		c.Send <- payload
	}
}

// mustMarshal is a helper function that marshals a value to JSON and logs any errors, returning an empty JSON object if marshalling fails.
func (r *Room) mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		log.Println("marshal error:", err)
		return []byte("{}")
	}
	return b
}
