// internal/game/game_logic.go
// This file contains the core game logic for Skribble, including turn management, scoring, and state transitions.
// It defines the necessary methods to handle player actions, update game state, and manage the flow of the game.
package game

import "errors"

// Start initializes the game and starts the first turn.
func (g *Game) Start() error {
	if len(g.Players) < 2 {
		return errors.New("not enough players")
	}

	if g.State == Playing {
		return nil
	}

	g.State = Playing
	g.playerIndex = 0
	return g.startNextTurn()
}
