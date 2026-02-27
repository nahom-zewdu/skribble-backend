// internal/engine/engine.go
// This file defines the Engine struct, which is responsible for managing the game loop and coordinating game state updates in the Skribble backend.
// The Engine periodically checks for turn timeouts and triggers appropriate game events to ensure smooth gameplay.

package engine

import (
	"time"

	"github.com/nahom-zewdu/skribble-backend/internal/game"
)

type Engine struct {
	game         *game.Game
	tickInterval time.Duration
}
