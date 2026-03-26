// internal/engine/engine.go
// Engine manages the game loop and coordinates game state updates in Skribble.
// It ensures timeouts and game events are processed safely, respecting domain invariants.

package engine

import (
	"github.com/nahom-zewdu/skribble-backend/internal/game"
)

type Engine struct {
	game *game.Game
}

func New(g *game.Game) *Engine {
	return &Engine{
		game: g,
	}
}

// Tick triggers game timeouts and returns any resulting events.
func (e *Engine) Tick() []game.GameEvent {
	events, _ := e.game.HandleTimeouts()
	return events
}

// Start initializes the game via the game domain.
func (e *Engine) Start() ([]game.GameEvent, error) {
	return e.game.Start()
}

// Guess processes a player's guess.
func (e *Engine) Guess(playerID, guess string) ([]game.GameEvent, error) {
	return e.game.Guess(playerID, guess)
}

// SelectWord allows the drawer to select a word.
func (e *Engine) SelectWord(playerID, word string) ([]game.GameEvent, error) {
	return e.game.SelectWord(playerID, word)
}

// AutoSelectWord triggers automatic word selection if deadline exceeded.
func (e *Engine) AutoSelectWord() ([]game.GameEvent, error) {
	return e.game.AutoSelectWord()
}

// EndTurn ends the current turn safely.
func (e *Engine) EndTurn() ([]game.GameEvent, error) {
	return e.game.EndTurn()
}

// AddPlayer adds a player through the game domain.
func (e *Engine) AddPlayer(id, name string) ([]game.GameEvent, error) {
	events, err := e.game.AddPlayer(id, name)
	if err != nil {
		return nil, err
	}

	startEvents, err := e.game.Start()
	if err == nil {
		events = append(events, startEvents...)
	}

	return events, nil
}

// RemovePlayer removes a player through the game domain.
func (e *Engine) RemovePlayer(id string) ([]game.GameEvent, error) {
	return e.game.RemovePlayer(id)
}

// Snapshot returns a read-only snapshot of the current game state.
func (e *Engine) Snapshot() game.GameSnapshot {
	return e.game.Snapshot()
}

// CurrentDrawerID returns the ID of the current drawer, or empty string if no active turn.
func (e *Engine) CurrentDrawerID() string {
	if e.game.CurrentTurn == nil {
		return ""
	}
	return e.game.CurrentTurn.DrawerID
}
