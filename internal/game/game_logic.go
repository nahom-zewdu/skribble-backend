// internal/game/game_logic.go
// This file contains the core game logic for Skribble, including turn management, scoring, and state transitions.
// It defines the necessary methods to handle player actions, update game state, and manage the flow of the game.
package game

import (
	"errors"
	"time"
)

/*
Domain Invariants Enforced Here:

1. Game state transitions only happen inside domain methods.
2. Turn lifecycle is fully controlled by Game.
3. Timeouts are handled ONLY through HandleTimeouts().
4. No external mutation of Game.State.
5. All meaningful transitions emit events.
*/

// Start initializes the game and starts the first turn.
func (g *Game) Start() ([]GameEvent, error) {
	if len(g.Players) < 2 {
		return nil, errors.New("not enough players")
	}

	if g.State == Playing {
		return nil, nil
	}

	g.State = Playing
	g.playerIndex = 0

	events := []GameEvent{
		{
			Type:      EventGameStarted,
			Timestamp: time.Now(),
		},
	}

	turnEvents, err := g.startNextTurn()
	if err != nil {
		return events, err
	}

	return append(events, turnEvents...), nil
}

// startNextTurn sets up the next turn in the game, selecting the next drawer and generating word choices.
func (g *Game) startNextTurn() ([]GameEvent, error) {
	if g.CurrentTurn != nil && g.CurrentTurn.Number >= g.MaxTurns {
		g.State = Ended

		return []GameEvent{
			{
				Type:      EventGameEnded,
				Timestamp: time.Now(),
				Payload: GameEndedPayload{
					Players: g.Players,
				},
			},
		}, errors.New("game ended")
	}

	if len(g.Players) < 2 {
		g.State = Waiting
		return nil, errors.New("not enough players to continue")
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
		Choices:           g.wordProvider.GenerateChoices(3),
		Phase:             PhaseSelecting,
		SelectionDeadline: now.Add(10 * time.Second),
		Guessed:           make(map[string]bool),
		Completed:         false,
	}

	g.playerIndex++

	return []GameEvent{
		{
			Type:      EventTurnStarted,
			Timestamp: now,
			Payload: TurnStartedPayload{
				TurnNumber: turnNumber,
				DrawerID:   drawer.ID,
				Choices:    g.CurrentTurn.Choices,
			},
		},
		{
			Type:      EventWordSelectionStart,
			Timestamp: now,
		},
	}, nil
}

// SelectWord allows the drawer to select a word for the current turn, transitioning to the drawing phase.
func (g *Game) SelectWord(playerID, word string) ([]GameEvent, error) {
	if g.CurrentTurn == nil {
		return nil, errors.New("no active turn")
	}

	if g.CurrentTurn.Phase != PhaseSelecting {
		return nil, errors.New("not in selection phase")
	}

	if playerID != g.CurrentTurn.DrawerID {
		return nil, errors.New("only drawer can select word")
	}

	valid := false
	for _, w := range g.CurrentTurn.Choices {
		if w == word {
			valid = true
			break
		}
	}
	if !valid {
		return nil, errors.New("invalid word choice")
	}

	now := time.Now()

	g.CurrentTurn.Word = word
	g.CurrentTurn.Phase = PhaseDrawing
	g.CurrentTurn.StartTime = now
	g.CurrentTurn.PlayDeadline = now.Add(65 * time.Second)

	return []GameEvent{
		{
			Type:      EventWordSelected,
			Timestamp: now,
			Payload: WordSelectedPayload{
				DrawerID: playerID,
				Word:     word,
			},
		},
	}, nil
}

// AutoSelectWord automatically selects a word for the drawer if they fail to choose within the deadline.
func (g *Game) AutoSelectWord() ([]GameEvent, error) {
	if g.CurrentTurn == nil {
		return nil, errors.New("no active turn")
	}

	if g.CurrentTurn.Phase != PhaseSelecting {
		return nil, errors.New("not in selection phase")
	}

	if time.Now().After(g.CurrentTurn.SelectionDeadline) {
		// fallback to first word
		word := g.CurrentTurn.Choices[0]
		return g.SelectWord(g.CurrentTurn.DrawerID, word)
	}
	return nil, errors.New("couldn't select word")
}

// Guess processes a player's guess and updates scores if correct.
func (g *Game) Guess(playerID, guess string) ([]GameEvent, error) {
	if g.CurrentTurn == nil {
		return nil, errors.New("no active turn")
	}

	if g.CurrentTurn.Completed {
		return nil, errors.New("turn is completed")
	}

	if playerID == g.CurrentTurn.DrawerID {
		return nil, errors.New("drawer cannot guess")
	}

	if guess != g.CurrentTurn.Word {
		return nil, nil // incorrect guess, no events
	}

	if g.CurrentTurn.Guessed[playerID] {
		return nil, errors.New("already guessed")
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

	return []GameEvent{
		{
			Type:      EventCorrectGuess,
			Timestamp: time.Now(),
			Payload: CorrectGuessPayload{
				PlayerID: playerID,
				Score:    score,
			},
		},
	}, nil
}

// EndTurn marks the current turn as completed and starts the next turn.
func (g *Game) EndTurn() ([]GameEvent, error) {
	if g.CurrentTurn == nil {
		return nil, errors.New("no active turn")
	}

	g.CurrentTurn.Completed = true

	events := []GameEvent{
		{
			Type:      EventTurnEnded,
			Timestamp: time.Now(),
			Payload: TurnEndedPayload{
				TurnNumber: g.CurrentTurn.Number,
				Word:       g.CurrentTurn.Word,
			},
		},
	}

	nextEvents, err := g.startNextTurn()
	events = append(events, nextEvents...)

	return events, err
}

// CheckTurnTimeout checks if the current turn has exceeded its play deadline and ends the turn if necessary.
func (g *Game) CheckTurnTimeout() ([]GameEvent, error) {
	if g.CurrentTurn == nil {
		return nil, errors.New("no active turn")
	}

	if g.CurrentTurn.Phase == PhaseDrawing &&
		time.Now().After(g.CurrentTurn.PlayDeadline) {

		return g.EndTurn()
	}
	return nil, nil
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
