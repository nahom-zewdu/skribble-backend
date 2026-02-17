// internal/room/manager.go
// This file defines the Manager struct, which is responsible for managing multiple game rooms in the Skribble backend. It provides methods to create, retrieve, and delete rooms.

package room

import (
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

func (m *Manager) CreateRoom(id string) *Room {
	m.mu.Lock()
	defer m.mu.Unlock()

	room := NewRoom(id)
	m.rooms[id] = room
	return room
}

func (m *Manager) GetRoom(id string) (*Room, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	room, ok := m.rooms[id]
	return room, ok
}

func (m *Manager) DeleteRoom(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.rooms, id)
}
