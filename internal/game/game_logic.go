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

// SetWord assigns the selected word to the current turn and starts the timer.
func (g *Game) SetWord(word string) error {
	if g.CurrentTurn == nil {
		return errors.New("no active turn")
	}

	g.CurrentTurn.Word = word
	g.CurrentTurn.StartTime = time.Now()
	return nil
}

// Guess processes a player's guess and updates scores if correct.
func (g *Game) Guess(playerID, guess string) (bool, error) {
	if g.CurrentTurn == nil {
		return false, errors.New("no active turn")
	}

	if g.CurrentTurn.Completed {
		return false, nil
	}

	if playerID == g.CurrentTurn.DrawerID {
		return false, nil
	}

	if guess != g.CurrentTurn.Word {
		return false, nil
	}

	// correct guess
	g.CurrentTurn.Guessed[playerID] = true

	elapsed := time.Since(g.CurrentTurn.StartTime).Seconds()
	score := int(100 - elapsed) // simple time-based scoring
	if score < 10 {
		score = 10
	}

	for _, p := range g.Players {
		if p.ID == playerID {
			p.Score += score
		}
	}

	return true, nil
}

// EndTurn marks the current turn as completed and starts the next turn.
func (g *Game) EndTurn() error {
	if g.CurrentTurn == nil {
		return errors.New("no active turn")
	}

	g.CurrentTurn.Completed = true
	return g.startNextTurn()
}
