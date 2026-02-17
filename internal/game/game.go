// internal/game/game.go
// This file defines the Game struct, which represents the state of a game session in the Skribble backend.
// It includes fields for the current turn, maximum turns, and game state, along with a constructor function to initialize a new game.

package game

type State string

const (
	Waiting State = "waiting"
	Playing State = "playing"
	Paused  State = "paused"
	Ended   State = "ended"
)

type Game struct {
	CurrentTurn int
	MaxTurns    int
	State       State
}

func NewGame() *Game {
	return &Game{
		CurrentTurn: 0,
		MaxTurns:    9, // fixed as per your rule
		State:       Waiting,
	}
}
