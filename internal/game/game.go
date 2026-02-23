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

type Player struct {
	ID    string
	Name  string
	Score int
}

type Turn struct {
	Number    int
	DrawerID  string
	Word      string
	StartTime time.Time
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
