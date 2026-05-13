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
	EventSelectionTimeout     EventType = "selection_timeout"
	EventDrawingTimeout       EventType = "drawing_timeout"
)

type GameEvent struct {
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"data"`
}

// --------------------
// Payloads
// --------------------

type TurnStartedPayload struct {
	TurnNumber int    `json:"turnNumber"`
	DrawerID   string `json:"drawerID"`
}

type WordSelectedPayload struct {
	DrawerID     string    `json:"drawerID"`
	Word         string    `json:"word"`
	PlayDeadline time.Time `json:"playDeadline"`
}

type CorrectGuessPayload struct {
	PlayerID   string `json:"playerID"`
	PlayerName string `json:"playerName"`

	Score int `json:"score"`

	DrawerID     string `json:"drawerID"`
	DrawerPoints int    `json:"drawerPoints"`

	TotalScore int `json:"totalScore"`
}

type TurnEndedPayload struct {
	TurnNumber        int              `json:"turnNumber"`
	Word              string           `json:"word"`
	Players           []PlayerSnapshot `json:"players"`
	NextTurnStartTime time.Time        `json:"nextTurnStartTime"`
}

type GameEndedPayload struct {
	Players     []PlayerSnapshot `json:"players"`
	RestartTime time.Time        `json:"restartTime"`
}

type PlayerJoinedPayload struct {
	PlayerID string `json:"playerID"`
	Name     string `json:"name"`
}

type PlayerLeftPayload struct {
	PlayerID string `json:"playerID"`
	Name     string `json:"name"`
}

type WordSelectionStartedPayload struct {
	DrawerID string    `json:"drawerID"`
	Choices  []string  `json:"choices"`
	Deadline time.Time `json:"deadline"`
}
