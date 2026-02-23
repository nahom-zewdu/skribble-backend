// internal/game/game_logic.go
// This file contains the core game logic for Skribble, including turn management, scoring, and state transitions.
// It defines the necessary methods to handle player actions, update game state, and manage the flow of the game.
package game

import (
	"errors"
	"time"
)

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

// startNextTurn advances the game to the next turn, selecting the next drawer and resetting turn state.
func (g *Game) startNextTurn() error {
	if g.CurrentTurn != nil && g.CurrentTurn.Number >= g.MaxTurns {
		g.State = Ended
		return errors.New("game ended")
	}

	if len(g.Players) == 0 {
		g.State = Waiting
		return errors.New("no players")
	}

	drawer := g.Players[g.playerIndex%len(g.Players)]

	turnNumber := 1
	if g.CurrentTurn != nil {
		turnNumber = g.CurrentTurn.Number + 1
	}

	g.CurrentTurn = &Turn{
		Number:    turnNumber,
		DrawerID:  drawer.ID,
		Word:      "", // word selected later
		StartTime: time.Now(),
		Guessed:   make(map[string]bool),
	}

	g.playerIndex++

	return nil
}
