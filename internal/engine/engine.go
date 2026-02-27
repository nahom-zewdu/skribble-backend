// internal/engine/engine.go
// This file defines the Engine struct, which is responsible for managing the game loop and coordinating game state updates in the Skribble backend.
// The Engine periodically checks for turn timeouts and triggers appropriate game events to ensure smooth gameplay.

package engine

import (
	"time"

	"github.com/nahom-zewdu/skribble-backend/internal/game"
)

type Engine struct {
	game         *game.Game
	tickInterval time.Duration
}

func New(g *game.Game) *Engine {
	return &Engine{
		game:         g,
		tickInterval: 250 * time.Millisecond,
	}
}

// Tick checks for turn timeouts and triggers game events accordingly.
// It should be called periodically (e.g., every 250ms) to ensure timely updates.
func (e *Engine) Tick() []game.GameEvent {
	var events []game.GameEvent

	// Word auto-select
	if e.game.CurrentTurn != nil &&
		e.game.CurrentTurn.Phase == game.PhaseSelecting &&
		time.Now().After(e.game.CurrentTurn.SelectionDeadline) {

		ev, _ := e.game.SelectWord(
			e.game.CurrentTurn.DrawerID,
			e.game.CurrentTurn.Choices[0],
		)
		events = append(events, ev...)
	}

	// Play timeout
	if e.game.CurrentTurn != nil &&
		e.game.CurrentTurn.Phase == game.PhaseDrawing &&
		time.Now().After(e.game.CurrentTurn.PlayDeadline) {

		ev, _ := e.game.EndTurn()
		events = append(events, ev...)
	}

	return events
}

// Start initiates the game and returns any resulting game events.
func (e *Engine) Start() ([]game.GameEvent, error) {
	return e.game.Start()
}

// Guess processes a player's guess and returns any resulting game events.
func (e *Engine) Guess(playerID, guess string) ([]game.GameEvent, error) {
	return e.game.Guess(playerID, guess)
}

// SelectWord allows the drawer to select a word and returns any resulting game events.
func (e *Engine) SelectWord(playerID, word string) ([]game.GameEvent, error) {
	return e.game.SelectWord(playerID, word)
}

// AutoSelectWord automatically selects a word for the drawer if they fail to choose within the deadline.
func (e *Engine) AutoSelectWord() ([]game.GameEvent, error) {
	return e.game.AutoSelectWord()
}

// EndTurn ends the current turn and returns any resulting game events.
func (e *Engine) EndTurn() ([]game.GameEvent, error) {
	return e.game.EndTurn()
}

// AddPlayer adds a new player to the game.
func (e *Engine) AddPlayer(id, name string) {
	e.game.AddPlayer(id, name)
}

// RemovePlayer removes a player from the game by their ID.
func (e *Engine) RemovePlayer(id string) {
	e.game.RemovePlayer(id)
}

// Game returns the current game state.
func (e *Engine) Game() *game.Game {
	return e.game
}

// UpdateGameState allows external updates to the game state, if needed.
func (e *Engine) UpdateGameState(newState game.State) {
	e.game.UpdateGameState(newState)
}
