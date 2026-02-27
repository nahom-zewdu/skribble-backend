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

// AddPlayer adds a new player to the game.
func (g *Game) AddPlayer(id, name string) {
	g.Players = append(g.Players, &Player{
		ID:   id,
		Name: name,
	})
}

// RemovePlayer removes a player and handles structural consequences.
// It emits domain events if the removal affects game flow.
func (g *Game) RemovePlayer(id string) ([]GameEvent, error) {
	if len(g.Players) == 0 {
		return
	}

	removedIndex := -1
	for i, p := range g.Players {
		if p.ID == id {
			removedIndex = i
			break
		}
	}

	if removedIndex == -1 {
		return
	}

	// Check if removed player is current drawer
	isDrawer := false
	if g.CurrentTurn != nil && g.CurrentTurn.DrawerID == id {
		isDrawer = true
	}

	// Remove player from slice
	g.Players = append(g.Players[:removedIndex], g.Players[removedIndex+1:]...)

	// Adjust rotation index safely
	if removedIndex < g.playerIndex {
		g.playerIndex--
	}

	if g.playerIndex >= len(g.Players) {
		g.playerIndex = 0
	}

	// If drawer left, mark turn as completed (engine will decide next action)
	if isDrawer && g.CurrentTurn != nil {
		g.CurrentTurn.Phase = PhaseEnded
		g.CurrentTurn.Completed = true
	}

	// If no players left → end game
	if len(g.Players) == 0 {
		g.State = Ended
		return
	}

	// If fewer than 2 players → back to waiting
	if len(g.Players) < 2 {
		g.State = Waiting
	}
}
