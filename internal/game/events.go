// internal/game/events.go
// This file defines the event system for Skribble, including the GameEvent struct and event types.
// The GameEvent struct is used to represent various events that occur during the game, such as game start, turn start, word selection, correct guesses, and game end.
// Each event includes a type, timestamp, and an optional payload for additional data related to the event.
package game

import "time"

type EventType string

// Event types for the game
const (
	EventGameStarted          EventType = "game_started"
	EventTurnStarted          EventType = "turn_started"
	EventWordSelected         EventType = "word_selected"
	EventCorrectGuess         EventType = "correct_guess"
	EventTurnEnded            EventType = "turn_ended"
	EventGameEnded            EventType = "game_ended"
	EventPlayerJoined         EventType = "player_joined"
	EventPlayerLeft           EventType = "player_left"
	EventWordSelectionStarted EventType = "word_selection_started"
)

type GameEvent struct {
	Type      EventType
	Timestamp time.Time
	Payload   interface{}
}

type TurnStartedPayload struct {
	TurnNumber int
	DrawerID   string
}

type WordSelectedPayload struct {
	DrawerID     string
	Word         string
	PlayDeadline time.Time `json:"playDeadline"`
}

type CorrectGuessPayload struct {
	PlayerID string
	Score    int
}

type TurnEndedPayload struct {
	TurnNumber        int
	Word              string
	Players           []PlayerSnapshot
	NextTurnStartTIme time.Time
}

type GameEndedPayload struct {
	Players     []PlayerSnapshot
	RestartTime time.Time `json:"restartTime"`
}

type PlayerJoinedPayload struct {
	PlayerID string
	Name     string
}

type PlayerLeftPayload struct {
	PlayerID string
	Name     string
}

type WordSelectionStartedPayload struct {
	DrawerID string
	Choices  []string
	Deadline time.Time
}
