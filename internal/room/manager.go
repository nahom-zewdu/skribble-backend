// internal/room/manager.go
// This file defines the Manager struct, which is responsible for managing multiple game rooms in the Skribble backend.
package room

import (
	"log"
	"sync"
)

type Manager struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
	}
}

func (m *Manager) GetOrCreateRoom(id string) *Room {
	m.mu.Lock()
	defer m.mu.Unlock()

	if room, ok := m.rooms[id]; ok {
		return room
	}

	room := NewRoom(id, m)
	m.rooms[id] = room
	return room
}

// DeleteRoom removes a room from the manager. It acquires a write lock to safely modify the rooms map, deletes the specified room, and logs the deletion.
func (m *Manager) DeleteRoom(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.rooms, id)
	log.Printf("[room manager] Deleted room %s\n", id)
}

// ROOM CODE

func generateRoomCode() string {
	length := 6
	b := make([]byte, length)

	for i := range b {
		b[i] = roomCodeCharset[rand.Intn(len(roomCodeCharset))]
	}

	return string(b)
}