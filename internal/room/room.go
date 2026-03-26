// internal/room/room.go
// This file defines the Room struct, which manages a single game room in the Skribble backend.
// It handles client registration, message processing, and game state management for the room.
// The Room struct interacts with the Game logic to manage turns, scoring, and player actions.
package room

import (
	"encoding/json"
	"log"
	"time"

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

// run is the main loop for the room.
// It delegates all game-related decisions to the Engine and only reacts to emitted events.
func (r *Room) run() {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:
			events := r.engine.Tick()
			r.handleEvents(events)

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

// onJoin handles new player registration and delegates domain updates to the Engine.
func (r *Room) onJoin(c *client.Client) {
	events, err := r.engine.AddPlayer(c.ID, c.Name)
	if err != nil {
		log.Println("engine AddPlayer error:", err)
		return
	}

	// Send full snapshot to the newly joined player
	snapshot := r.engine.Snapshot()
	snapshot.SelfID = c.ID // Include SelfID in the snapshot for client-side logic
	c.Send <- r.mustMarshal(transport.Message{
		Type: "game_snapshot",
		Data: snapshot,
	})

	r.handleEvents(events)
}

// onLeave handles player removal and delegates domain updates to the Engine.
func (r *Room) onLeave(c *client.Client) {
	events, err := r.engine.RemovePlayer(c.ID)
	if err != nil {
		log.Println("engine RemovePlayer error:", err)
		return
	}

	r.handleEvents(events)
}

// onMessage routes incoming client messages to the appropriate Engine command.
func (r *Room) onMessage(sender *client.Client, raw []byte) {
	var incoming transport.ClientMessage

	if err := json.Unmarshal(raw, &incoming); err != nil {
		return
	}

	switch incoming.Type {

	case "chat":
		r.handleChat(sender, incoming.Data)

	case "select_word":
		var data struct {
			Word string `json:"word"`
		}
		if err := json.Unmarshal(incoming.Data, &data); err != nil {
			return
		}

		events, err := r.engine.SelectWord(sender.ID, data.Word)
		if err != nil {
			log.Println("SelectWord error:", err)
			return
		}

		r.handleEvents(events)

	case "draw_start":
		r.handleDrawStart(sender, incoming.Data)

	case "draw_move":
		r.handleDrawMove(sender, incoming.Data)

	case "draw_end":
		r.handleDrawEnd(sender)

	case "clear_canvas":
		r.handleClearCanvas(sender)
	}
}

// handleChat sends guesses to the Engine and falls back to normal chat broadcast if no domain event is emitted.
func (r *Room) handleChat(sender *client.Client, raw json.RawMessage) {
	var chat transport.ChatMessage
	if err := json.Unmarshal(raw, &chat); err != nil {
		return
	}

	events, err := r.engine.Guess(sender.ID, chat.Text)
	if err != nil {
		log.Println("Guess error:", err)
		return
	}

	if len(events) > 0 {
		r.handleEvents(events)
		return
	}

	// Otherwise it's normal chat
	r.broadcastChat(sender.Name, chat.Text)
}

// handleDrawStart processes the start of a drawing action from the drawer and broadcasts it to other clients.
func (r *Room) handleDrawStart(sender *client.Client, raw json.RawMessage) {
	if !r.isDrawer(sender.ID) {
		return
	}

	var data transport.DrawStart
	if err := json.Unmarshal(raw, &data); err != nil {
		return
	}

	r.broadcastDraw("draw_start", map[string]interface{}{
		"x":        data.X,
		"y":        data.Y,
		"senderID": sender.ID,
	})
}

// handleDrawMove processes drawing movements from the drawer and broadcasts them to other clients.
func (r *Room) handleDrawMove(sender *client.Client, raw json.RawMessage) {
	if !r.isDrawer(sender.ID) {
		return
	}

	var data transport.DrawMove
	if err := json.Unmarshal(raw, &data); err != nil {
		return
	}

	// ⚠️ Guard against abuse (VERY important)
	if len(data.Points) > 100 {
		return
	}

	r.broadcastDraw("draw_move", map[string]interface{}{
		"points":   data.Points,
		"senderID": sender.ID,
	})
}

// handleDrawEnd processes the end of a drawing action and broadcasts it to other clients.
func (r *Room) handleDrawEnd(sender *client.Client) {
	if !r.isDrawer(sender.ID) {
		return
	}

	r.broadcastDraw("draw_end", map[string]interface{}{
		"senderID": sender.ID,
	})
}

// handleEvents translates domain GameEvents into transport-level messages.
func (r *Room) handleEvents(events []game.GameEvent) {
	for _, e := range events {

		switch e.Type {

		case game.EventPlayerJoined:
			payload := e.Payload.(game.PlayerJoinedPayload)
			r.broadcastSystem(payload.Name + " joined")

		case game.EventPlayerLeft:
			payload := e.Payload.(game.PlayerLeftPayload)
			r.broadcastSystem(payload.Name + " left")

		case game.EventGameStarted:
			r.broadcastSystem("Game started")

		case game.EventTurnStarted:
			payload := e.Payload.(game.TurnStartedPayload)

			for _, c := range r.clients {
				c.Send <- r.mustMarshal(transport.Message{
					Type: "turn_started",
					Data: payload,
				})
			}

		case game.EventWordSelectionStarted:
			payload := e.Payload.(game.WordSelectionStartedPayload)

			for _, c := range r.clients {
				c.Send <- r.mustMarshal(transport.Message{
					Type: "word_selection_started",
					Data: payload,
				})
			}

		case game.EventWordSelected:
			payload := e.Payload.(game.WordSelectedPayload)

			for _, c := range r.clients {
				if c.ID == payload.DrawerID {
					c.Send <- r.mustMarshal(transport.Message{
						Type: "drawing_started",
						Data: map[string]interface{}{
							"word":     payload.Word,
							"deadline": payload.PlayDeadline,
						},
					})
				} else {
					c.Send <- r.mustMarshal(transport.Message{
						Type: "drawing_started",
						Data: map[string]interface{}{
							"deadline":   payload.PlayDeadline,
							"maskedWord": payload.Word, // Note: this will be masked (e.g. "_ _ _ _")
						},
					})
				}
			}

		case game.EventCorrectGuess:
			payload := e.Payload.(game.CorrectGuessPayload)

			r.broadcastCorrectGuess(payload)

		case game.EventTurnEnded:
			payload := e.Payload.(game.TurnEndedPayload)

			msg := transport.Message{
				Type: "turn_ended",
				Data: payload,
			}

			data := r.mustMarshal(msg)

			for _, c := range r.clients {
				c.Send <- data
			}

		case game.EventGameEnded:
			payload := e.Payload.(game.GameEndedPayload)

			msg := transport.Message{
				Type: "game_ended",
				Data: payload,
			}

			data := r.mustMarshal(msg)

			for _, c := range r.clients {
				c.Send <- data
			}
		}
	}
}

// broadcastChat sends a chat message to all clients in the room.
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

// broadcastCorrectGuess notifies clients that a player guessed correctly.
func (r *Room) broadcastCorrectGuess(guessPayload game.CorrectGuessPayload) {
	msg := transport.Message{
		Type: "correct_guess",
		Data: guessPayload,
	}

	payload := r.mustMarshal(msg)

	for _, c := range r.clients {
		c.Send <- payload
	}
}

// broadcastSystem sends a system message to all clients.
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

// mustMarshal marshals a value to JSON and logs errors if any.
func (r *Room) mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		log.Println("marshal error:", err)
		return []byte("{}")
	}
	return b
}
