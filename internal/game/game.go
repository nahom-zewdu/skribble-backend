// internal/game/game.go
// This file defines the core game logic for Skribble, including the game state machine, player management, and turn handling.
// It provides the necessary structures and methods to manage the game flow, track player scores, and handle game state transitions.
package game

import (
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
}

func NewGame() *Game {
	return &Game{
		State:    Waiting,
		MaxTurns: 9,
		Players:  []*Player{},
	}
}

// AddPlayer adds a new player to the game.
func (g *Game) AddPlayer(id, name string) {
	g.Players = append(g.Players, &Player{
		ID:   id,
		Name: name,
	})
}

// RemovePlayer removes a player from the game by their ID.
func (g *Game) RemovePlayer(id string) {
	newPlayers := []*Player{}
	for _, p := range g.Players {
		if p.ID != id {
			newPlayers = append(newPlayers, p)
		}
	}
	g.Players = newPlayers
}
