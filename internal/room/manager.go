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

func (m *Manager) GetPrivateRoom(id string) (*Room, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	room, ok := m.rooms[id]

	if !ok {
		return nil, false
	}

	if room.Type != PrivateRoom {
		return nil, false
	}

	return room, true
}

// DELETE

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