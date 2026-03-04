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

	if !g.CanStart() {
		return nil, errors.New("game cannot start yet")
	}

	g.Reset() // Ensure game is reset before starting

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
	if len(g.Players) < 2 {
		g.State = Waiting
		return nil, errors.New("not enough players to continue")
	}

	// If max turns reached, end game
	if g.CurrentTurn != nil && g.CurrentTurn.Number >= g.MaxTurns {
		restart := time.Now().Add(5 * time.Second)
		g.RestartDeadline = &restart
		g.State = Ended
		return []GameEvent{
			{
				Type:      EventGameEnded,
				Timestamp: time.Now(),
				Payload: GameEndedPayload{
					Players:     g.Players,
					RestartTime: restart,
				},
			},
		}, errors.New("game ended")
	}

	// Select next drawer in round-robin fashion
	drawer := g.Players[g.playerIndex%len(g.Players)]
	turnNumber := 1
	if g.CurrentTurn != nil {
		turnNumber = g.CurrentTurn.Number + 1
	}

	// Initialize new turn
	now := time.Now()
	choices := g.wordProvider.GenerateChoices(3)
	g.CurrentTurn = &Turn{
		Number:            turnNumber,
		DrawerID:          drawer.ID,
		Choices:           choices,
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
			},
		},
		{
			Type:      EventWordSelectionStarted,
			Timestamp: now,
			Payload: WordSelectionStartedPayload{
				DrawerID: drawer.ID,
				Choices:  choices,
				Deadline: g.CurrentTurn.SelectionDeadline,
				// Note: Deadline is included here for clients to show selection timer, but actual enforcement is in HandleTimeouts()
				// This ensures that even if a client misses the deadline, the game will skip to the next turn without relying on client-side timers.
			},
		},
	}, nil
}

// SelectWord allows the drawer to select a word for the current turn.
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
				DrawerID:     playerID,
				Word:         word,
				PlayDeadline: g.CurrentTurn.PlayDeadline,
			},
		},
	}, nil
}

// AutoSelectWord automatically chooses a word if drawer misses the deadline.
func (g *Game) AutoSelectWord() ([]GameEvent, error) {
	if g.CurrentTurn == nil {
		return nil, errors.New("no active turn")
	}
	if g.CurrentTurn.Phase != PhaseSelecting {
		return nil, errors.New("not in selection phase")
	}
	if time.Now().After(g.CurrentTurn.SelectionDeadline) {
		return g.SelectWord(g.CurrentTurn.DrawerID, g.CurrentTurn.Choices[0])
	}
	return nil, nil
}

// Guess processes a player's guess.
func (g *Game) Guess(playerID, guess string) ([]GameEvent, error) {
	if g.CurrentTurn == nil || g.CurrentTurn.Completed {
		return nil, errors.New("no active turn or turn completed")
	}
	if playerID == g.CurrentTurn.DrawerID {
		return nil, errors.New("drawer cannot guess")
	}
	if g.CurrentTurn.Guessed[playerID] {
		return nil, errors.New("already guessed")
	}
	if guess != g.CurrentTurn.Word {
		return nil, nil
	}

	g.CurrentTurn.Guessed[playerID] = true

	// Time-based scoring
	elapsed := time.Since(g.CurrentTurn.StartTime).Seconds()
	maxDuration := 65.0
	remaining := maxDuration - elapsed
	if remaining < 0 {
		remaining = 0
	}
	score := int((remaining / maxDuration) * 100)
	if score < 10 {
		score = 10
	}

	for _, p := range g.Players {
		if p.ID == playerID {
			p.Score += score
			break
		}
	}
	// Drawer gets half points
	for _, p := range g.Players {
		if p.ID == g.CurrentTurn.DrawerID {
			p.Score += score / 2
			break
		}
	}

	event := []GameEvent{
		{
			Type:      EventCorrectGuess,
			Timestamp: time.Now(),
			Payload: CorrectGuessPayload{
				PlayerID: playerID,
				Score:    score,
			},
		},
	}

	// Check if all guessers have finished
	if g.AllGuessersFinished() {
		endEvents, err := g.EndTurn()
		if err == nil {
			event = append(event, endEvents...)
		}
	}

	return event, nil
}

// EndTurn marks the current turn as completed and starts next turn or ends game.
func (g *Game) EndTurn() ([]GameEvent, error) {
	if g.CurrentTurn == nil || g.CurrentTurn.Completed {
		return nil, errors.New("no active turn or already ended")
	}

	g.CurrentTurn.Completed = true
	g.CurrentTurn.Phase = PhaseEnded

	now := time.Now()
	transitionDeadline := now.Add(3 * time.Second)
	events := []GameEvent{
		{
			Type:      EventTurnEnded,
			Timestamp: now,
			Payload: TurnEndedPayload{
				TurnNumber:        g.CurrentTurn.Number,
				Word:              g.CurrentTurn.Word,
				Players:           g.copyPlayers(),
				NextTurnStartTIme: transitionDeadline,
			},
		},
	}

	// If less than 2 players remain end game
	if len(g.Players) < 2 {
		restart := time.Now().Add(5 * time.Second)
		g.RestartDeadline = &restart
		g.State = Ended
		events = append(events, GameEvent{
			Type:      EventGameEnded,
			Timestamp: now,
			Payload: GameEndedPayload{
				Players:     g.copyPlayers(),
				RestartTime: restart,
			},
		})
		return events, nil
	}

	// Schedule next turn transition (3 seconds)
	transition := now.Add(3 * time.Second)
	g.TurnTransitionDeadline = &transition

	return events, nil
}

// HandleTimeouts processes word selection and drawing timeouts.
func (g *Game) HandleTimeouts() ([]GameEvent, error) {
	if g.CurrentTurn == nil || g.CurrentTurn.Completed {
		return nil, nil
	}

	now := time.Now()

	// Word selection timeout
	if g.CurrentTurn.Phase == PhaseSelecting && now.After(g.CurrentTurn.SelectionDeadline) {
		g.CurrentTurn.Word = "" // no word selected
		return g.EndTurn()
	}

	// Drawing timeout
	if g.CurrentTurn.Phase == PhaseDrawing && now.After(g.CurrentTurn.PlayDeadline) {
		return g.EndTurn()
	}

	// Handle automatic game restart
	if g.State == Ended && g.RestartDeadline != nil && g.CurrentTurn.Number == g.MaxTurns {
		if time.Now().After(*g.RestartDeadline) {
			g.RestartDeadline = nil
			g.Reset()
			return g.Start()
		}
	}

	// Handle turn transition after turn end
	if g.State == Playing &&
		g.TurnTransitionDeadline != nil &&
		time.Now().After(*g.TurnTransitionDeadline) {

		g.TurnTransitionDeadline = nil
		return g.startNextTurn()
	}

	return nil, nil
}

// MaskedWord returns the current word masked for guessing players.
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

// AllGuessersFinished checks if all guessers have guessed correctly.
func (g *Game) AllGuessersFinished() bool {
	if g.CurrentTurn == nil {
		return false
	}
	for _, p := range g.Players {
		if p.ID != g.CurrentTurn.DrawerID && !g.CurrentTurn.Guessed[p.ID] {
			return false
		}
	}
	return true
}

// Reset clears the game state for a new game while keeping players.
func (g *Game) Reset() {
	g.State = Waiting
	g.CurrentTurn = nil
	for _, p := range g.Players {
		p.Score = 0
	}
}

// copyPlayers creates a deep copy of the players slice to prevent external mutation.
func (g *Game) copyPlayers() []*Player {
	copied := make([]*Player, len(g.Players))
	for i, p := range g.Players {
		cp := *p
		copied[i] = &cp
	}
	return copied
}

// PlayerSnapshot returns a read-only snapshot of the current players for transport.
func (g *Game) PlayerSnapshot() []PlayerSnapshot {
	snapshot := make([]PlayerSnapshot, len(g.Players))
	for i, p := range g.Players {
		snapshot[i] = PlayerSnapshot{
			ID:    p.ID,
			Name:  p.Name,
			Score: p.Score,
		}
	}
	return snapshot
}

// Snapshot returns a read-only snapshot of the current game state for transport.
func (g *Game) Snapshot() GameSnapshot {
	snap := GameSnapshot{
		State:    g.State,
		MaxTurns: g.MaxTurns,
		Players:  g.PlayerSnapshot(),
	}

	if g.CurrentTurn != nil {
		snap.TurnNumber = g.CurrentTurn.Number
		snap.DrawerID = g.CurrentTurn.DrawerID
		snap.Phase = g.CurrentTurn.Phase
		snap.MaskedWord = g.MaskedWord()

		if !g.CurrentTurn.SelectionDeadline.IsZero() {
			d := g.CurrentTurn.SelectionDeadline
			snap.SelectionDeadline = &d
		}

		if !g.CurrentTurn.PlayDeadline.IsZero() {
			d := g.CurrentTurn.PlayDeadline
			snap.PlayDeadline = &d
		}
	}

	snap.TransitionDeadline = g.TurnTransitionDeadline
	snap.RestartDeadline = g.RestartDeadline

	return snap
}
