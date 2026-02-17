// internal/player/player.go
// This file defines the Player struct, which represents a player in the Skribble game. It includes fields for the player's ID, name, WebSocket connection, and score.

package player

import "github.com/gorilla/websocket"

type Player struct {
	ID   string
	Name string

	Conn  *websocket.Conn
	Score int
}
