// internal/room/manager.go
// This file defines the Manager struct responsible for managing game rooms in the Skribble backend.
// The Manager handles creating, retrieving, and deleting rooms, as well as generating unique room codes for public and private rooms.

package room

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/nahom-zewdu/skribble-backend/internal/metrics"
)

const roomCodeCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

type Manager struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

func NewManager() *Manager {
	rand.Seed(time.Now().UnixNano())

	return &Manager{
		rooms: make(map[string]*Room),
	}
}

// PUBLIC MATCHMAKING
func (m *Manager) FindOrCreatePublicRoom() *Room {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, room := range m.rooms {
		if room.Type != PublicRoom {
			continue
		}

		if room.IsJoinable() {
			return room
		}
	}

	roomID := generateRoomCode()

	room := NewRoom(roomID, PublicRoom, m)

	metrics.IncRooms() // Increment the active rooms metric when a new room is created
	m.rooms[roomID] = room

	log.Printf("[room manager] Created public room %s\n", roomID)

	return room
}

// PRIVATE ROOMS
func (m *Manager) CreatePrivateRoom() *Room {
	m.mu.Lock()
	defer m.mu.Unlock()

	roomID := generateRoomCode()

	room := NewRoom(roomID, PrivateRoom, m)

	metrics.IncRooms() // Increment the active rooms metric when a new room is created
	m.rooms[roomID] = room

	log.Printf("[room manager] Created private room %s\n", roomID)

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

// JOIN BY ID

func (m *Manager) GetRoom(id string) (*Room, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	room, ok := m.rooms[id]
	if !ok {
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
