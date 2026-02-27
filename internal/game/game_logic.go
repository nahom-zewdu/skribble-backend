// internal/game/game_logic.go
// This file contains the core game logic for Skribble, including turn management, scoring, and state transitions.
// It defines the necessary methods to handle player actions, update game state, and manage the flow of the game.
package game

import (
	"errors"
	"log"
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

	now := time.Now()

	g.CurrentTurn = &Turn{
		Number:            turnNumber,
		DrawerID:          drawer.ID,
		Choices:           generateWordChoices(),
		Phase:             PhaseSelecting,
		SelectionDeadline: now.Add(10 * time.Second),
		Guessed:           make(map[string]bool),
	}

	g.playerIndex++

	return nil
}

// SelectWord allows the drawer to select a word for the current turn, transitioning to the drawing phase.
func (g *Game) SelectWord(playerID, word string) error {
	if g.CurrentTurn == nil {
		return errors.New("no active turn")
	}

	if g.CurrentTurn.Phase != PhaseSelecting {
		return errors.New("not in selection phase")
	}

	if playerID != g.CurrentTurn.DrawerID {
		return errors.New("only drawer can select word")
	}

	// validate choice
	valid := false
	for _, w := range g.CurrentTurn.Choices {
		if w == word {
			valid = true
			break
		}
	}

	if !valid {
		return errors.New("invalid word choice")
	}

	now := time.Now()

	g.CurrentTurn.Word = word
	g.CurrentTurn.Phase = PhaseDrawing
	g.CurrentTurn.StartTime = now
	g.CurrentTurn.PlayDeadline = now.Add(65 * time.Second)

	return nil
}

// AutoSelectWord automatically selects a word for the drawer if they fail to choose within the deadline.
func (g *Game) AutoSelectWord() {
	if g.CurrentTurn == nil {
		return
	}

	if g.CurrentTurn.Phase != PhaseSelecting {
		return
	}

	if time.Now().After(g.CurrentTurn.SelectionDeadline) {
		// fallback to first word
		word := g.CurrentTurn.Choices[0]
		_ = g.SelectWord(g.CurrentTurn.DrawerID, word)
	}
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

	if g.CurrentTurn.Guessed[playerID] {
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

// CheckTurnTimeout checks if the current turn has exceeded its play deadline and ends the turn if necessary.
func (g *Game) CheckTurnTimeout() {
	if g.CurrentTurn == nil {
		return
	}

	if g.CurrentTurn.Phase == PhaseDrawing &&
		time.Now().After(g.CurrentTurn.PlayDeadline) {

		err := g.EndTurn()
		if err != nil {
			log.Printf("Error ending turn: %v", err)
		}
	}
}

// MaskedWord returns the current word with letters masked for guessing players.
func (g *Game) MaskedWord() string {
	if g.CurrentTurn == nil || g.CurrentTurn.Word == "" {
		return ""
	}

	mask := ""
	for range g.CurrentTurn.Word {
		mask += "_ "
	}
	return mask
}
