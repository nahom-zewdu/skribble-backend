// internal/room/room.go
// This file defines the Room struct, which manages a single game room in the Skribble backend.
// It handles client registration, message processing, and game state management for the room.
// The Room struct interacts with the Game logic to manage turns, scoring, and player actions.
package room

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/nahom-zewdu/skribble-backend/internal/client"
	"github.com/nahom-zewdu/skribble-backend/internal/engine"
	"github.com/nahom-zewdu/skribble-backend/internal/game"
	"github.com/nahom-zewdu/skribble-backend/internal/transport"
)

type RoomType string

const (
	PublicRoom  RoomType = "public"
	PrivateRoom RoomType = "private"
)

type Room struct {
	ID string

	Type RoomType

	MaxPlayers int

	clients map[string]*client.Client
	manager *Manager

	register   chan *client.Client
	unregister chan *client.Client
	incoming   chan clientMessage

	engine *engine.Engine
}

type clientMessage struct {
	client  *client.Client
	message []byte
}

func NewRoom(id string, roomType RoomType, m *Manager) *Room {
	g := game.NewGame()
	e := engine.New(g)

	r := &Room{
		ID:         id,
		Type:       roomType,
		clients:    make(map[string]*client.Client),
		MaxPlayers: 8,
		manager:    m,
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

func (r *Room) ClientCount() int {
	return len(r.clients)
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
			if r.ClientCount() <= 0 {
				log.Printf("[room] No clients left in room %s, deleting room\n", r.ID)
				r.manager.DeleteRoom(r.ID)

				return
			}

		case msg := <-r.incoming:
			r.onMessage(msg.client, msg.message)
		}
	}
}

// --------------------
// JOIN / LEAVE
// --------------------

func (r *Room) onJoin(c *client.Client) {
	events, err := r.engine.AddPlayer(c.ID, c.Name)
	if err != nil {
		log.Println("engine AddPlayer error:", err)
		return
	}

	snapshot := r.engine.Snapshot()
	snapshot.SelfID = c.ID

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

	log.Printf("incoming: %s from %s", incoming.Type, sender.ID)

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

// --------------------
// CHAT / GUESS
// --------------------

func (r *Room) handleChat(sender *client.Client, raw json.RawMessage) {
	var chat transport.ChatMessage
	if err := json.Unmarshal(raw, &chat); err != nil {
		return
	}

	// drawer or already-guessed players cannot chat
	if !r.engine.CanChat(sender.ID) {
		return
	}

	events, err := r.engine.Guess(sender.ID, chat.Text)
	if err != nil {
		log.Println("Guess error:", err)
	}

	if len(events) > 0 {
		r.handleEvents(events)
		return
	}

	r.broadcastChat(sender.Name, chat.Text)
}

// --------------------
// DRAWING HELPERS
// --------------------

func (r *Room) isDrawer(sender *client.Client) bool {
	snap := r.engine.Snapshot()
	return snap.DrawerID == sender.ID
}

func (r *Room) isDrawingPhase() bool {
	snap := r.engine.Snapshot()
	return snap.Phase == game.PhaseDrawing
}

// --------------------
// DRAW EVENTS
// --------------------

func (r *Room) handleDrawStart(sender *client.Client, raw json.RawMessage) {
	log.Println("draw_move accepted")

	if !r.isDrawer(sender) || !r.isDrawingPhase() {
		return
	}

	var data transport.DrawStart
	if err := json.Unmarshal(raw, &data); err != nil {
		return
	}

	r.broadcastDraw(sender.ID, "draw_start", data)
}

func (r *Room) handleDrawMove(sender *client.Client, raw json.RawMessage) {
	log.Println("draw_move accepted")

	if !r.isDrawer(sender) || !r.isDrawingPhase() {
		return
	}

	var data transport.DrawMove
	if err := json.Unmarshal(raw, &data); err != nil {
		return
	}

	// Abuse protection
	if len(data.Points) == 0 || len(data.Points) > 50 {
		return
	}

	for _, p := range data.Points {
		if p.X < 0 || p.Y < 0 || p.X > 2000 || p.Y > 2000 {
			return
		}
	}

	r.broadcastDraw(sender.ID, "draw_move", data)
}

// handleDrawEnd processes the end of a drawing action and broadcasts it to other clients.
func (r *Room) handleDrawEnd(sender *client.Client) {
	log.Println("draw_move accepted")

	if !r.isDrawer(sender) || !r.isDrawingPhase() {
		return
	}

	r.broadcastDraw(sender.ID, "draw_end", nil)
}

// handleClearCanvas processes a canvas clear action from the drawer and broadcasts it to other clients.
func (r *Room) handleClearCanvas(sender *client.Client) {
	log.Println("draw_move accepted")

	if !r.isDrawer(sender) || !r.isDrawingPhase() {
		return
	}

	r.broadcastAll("clear_canvas", nil)
}

// --------------------
// BROADCAST HELPERS
// --------------------

func (r *Room) broadcastDraw(senderID, eventType string, data interface{}) {
	msg := transport.Message{
		Type: eventType,
		Data: data,
	}

	payload := r.mustMarshal(msg)

	for _, c := range r.clients {
		if c.ID == senderID {
			continue
		}
		c.Send <- payload
	}
}

func (r *Room) broadcastAll(eventType string, data interface{}) {
	log.Println("broadcasting all:", eventType, data)

	msg := transport.Message{
		Type: eventType,
		Data: data,
	}

	payload := r.mustMarshal(msg)

	for _, c := range r.clients {
		c.Send <- payload
	}
}

// --------------------
// ENGINE EVENTS
// --------------------

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
			r.broadcastAll("clear_canvas", nil)

			payload := e.Payload.(game.TurnStartedPayload)
			r.broadcastAll("turn_started", payload)

		case game.EventWordSelectionStarted:
			payload := e.Payload.(game.WordSelectionStartedPayload)
			r.broadcastAll("word_selection_started", payload)

		case game.EventWordSelected:
			r.broadcastAll("clear_canvas", nil)

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
							"maskedWord": payload.Word,
						},
					})
				}
			}

		case game.EventCorrectGuess:
			payload := e.Payload.(game.CorrectGuessPayload)
			r.broadcastCorrectGuess(payload)

		case game.EventTurnEnded:
			payload := e.Payload.(game.TurnEndedPayload)
			r.broadcastAll("turn_ended", payload)

		case game.EventGameEnded:
			r.broadcastAll("clear_canvas", nil)

			payload := e.Payload.(game.GameEndedPayload)
			r.broadcastAll("game_ended", payload)

		case game.EventSelectionTimeout:
			payload := e.Payload.(map[string]interface{})
			r.broadcastAll("selection_timeout", payload)

		case game.EventDrawingTimeout:
			payload := e.Payload.(map[string]interface{})
			r.broadcastAll("drawing_timeout", payload)
		}
	}
}

// --------------------
// GENERIC BROADCASTS
// --------------------

func (r *Room) broadcastChat(sender, text string) {
	r.broadcastAll("chat", map[string]interface{}{
		"sender": sender,
		"text":   text,
	})
}

func (r *Room) broadcastCorrectGuess(payload game.CorrectGuessPayload) {
	// dedicated event
	r.broadcastAll("correct_guess", payload)

	// system chat message
	r.broadcastSystem(
		payload.PlayerName + " guessed correctly! +" + strconv.Itoa(payload.Score),
	)
}

func (r *Room) broadcastSystem(text string) {
	r.broadcastAll("system", transport.SystemMessage{Text: text})
}

// --------------------
// UTIL
// --------------------

func (r *Room) mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		log.Println("marshal error:", err)
		return []byte("{}")
	}
	return b
}

// IsFull checks if the room has reached its maximum player capacity by comparing the current client count with the MaxPlayers limit.
func (r *Room) IsFull() bool {
	return r.ClientCount() >= r.MaxPlayers
}

// IsJoinable determines if new players can join the room by checking if the room is not full.
func (r *Room) IsJoinable() bool {
	return !r.IsFull()
}
