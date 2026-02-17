// internal/room/room.go
// Manages game rooms, including player management and room state.
// Each room can have multiple players, and the room struct provides methods to add and remove players, as well as to get the current player count.
package room

import (
	"sync"

	"github.com/nahom-zewdu/skribble-backend/internal/player"
)

type Room struct {
	ID string

	Players map[string]*player.Player

	mu sync.RWMutex
}

func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		Players: make(map[string]*player.Player),
	}
}

func (r *Room) AddPlayer(p *player.Player) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Players[p.ID] = p
}

func (r *Room) RemovePlayer(playerID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.Players, playerID)
}

func (r *Room) PlayerCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.Players)
}
