// internal/game/game.go
// This file defines the core game logic for Skribble, including the game state machine, player management, and turn handling.
// It provides the necessary structures and methods to manage the game flow, track player scores, and handle game state transitions.
package game

import (
	"errors"
	"time"
)

/*
Game State Machine
*/

type State string

const (
	Waiting State = "waiting"
	Playing State = "playing"
	Paused  State = "paused"
	Ended   State = "ended"
)

type TurnPhase string

const (
	PhaseSelecting TurnPhase = "selecting"
	PhaseDrawing   TurnPhase = "drawing"
	PhaseEnded     TurnPhase = "ended"
)

type Player struct {
	ID    string
	Name  string
	Score int
}

type Turn struct {
	Number   int
	DrawerID string

	// Word system
	Word    string   // final selected word
	Choices []string // 3 selectable words
	Phase   TurnPhase

	// Timing
	SelectionDeadline time.Time
	PlayDeadline      time.Time
	StartTime         time.Time

	Guessed   map[string]bool
	Completed bool
}

type Game struct {
	State       State
	MaxTurns    int
	CurrentTurn *Turn

	Players     []*Player
	playerIndex int

	wordProvider WordProvider
}

func NewGame() *Game {
	return &Game{
		State:        Waiting,
		MaxTurns:     9,
		Players:      []*Player{},
		wordProvider: NewStaticWordProvider(),
	}
}

// AddPlayer adds a new player to the game and emits a PlayerJoined event.
func (g *Game) AddPlayer(id, name string) ([]GameEvent, error) {
	// Check if player already exists
	for _, p := range g.Players {
		if p.ID == id {
			return nil, errors.New("player already exists")
		}
	}

	player := &Player{
		ID:   id,
		Name: name,
	}
	var events []GameEvent

	g.Players = append(g.Players, player)

	events = append(events, GameEvent{
		Type:      EventPlayerJoined,
		Timestamp: time.Now(),
		Payload: PlayerJoinedPayload{
			PlayerID: player.ID,
			Name:     player.Name,
		},
	})

	return events, nil
}

// RemovePlayer removes a player and emits PlayerLeft and potentially GameEnded events.
func (g *Game) RemovePlayer(id string) ([]GameEvent, error) {
	if len(g.Players) == 0 {
		return nil, nil
	}

	removedIndex := -1
	var removedPlayer *Player
	for i, p := range g.Players {
		if p.ID == id {
			removedIndex = i
			removedPlayer = p
			break
		}
	}

	if removedIndex == -1 {
		return nil, nil
	}

	isDrawer := false
	if g.CurrentTurn != nil && g.CurrentTurn.DrawerID == id {
		isDrawer = true
	}

	// Remove from slice
	g.Players = append(g.Players[:removedIndex], g.Players[removedIndex+1:]...)

	// Fix rotation index
	if removedIndex < g.playerIndex {
		g.playerIndex--
	}
	if g.playerIndex >= len(g.Players) {
		g.playerIndex = 0
	}

	events := []GameEvent{
		{
			Type:      EventPlayerLeft,
			Timestamp: time.Now(),
			Payload: PlayerLeftPayload{
				PlayerID: removedPlayer.ID,
				Name:     removedPlayer.Name,
			},
		},
	}

	// If less than 2 players remain → end game
	if len(g.Players) < 2 {
		g.State = Ended
		if g.CurrentTurn != nil && !g.CurrentTurn.Completed {
			g.CurrentTurn.Completed = true
			g.CurrentTurn.Phase = PhaseEnded
		}

		events = append(events, GameEvent{
			Type:      EventGameEnded,
			Timestamp: time.Now(),
			Payload: GameEndedPayload{
				Players: g.Players,
			},
		})
		return events, nil
	}

	// If drawer left, end turn through domain
	if isDrawer && g.CurrentTurn != nil && !g.CurrentTurn.Completed {
		turnEvents, err := g.EndTurn()
		events = append(events, turnEvents...)
		return events, err
	}

	return events, nil
}

// CanStart checks if the game has enough players to start.
func (g *Game) CanStart() bool {
	return len(g.Players) >= 2
}
