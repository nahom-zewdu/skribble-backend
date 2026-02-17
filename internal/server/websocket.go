// internal/server/websocket.go
// This file defines the WebSocket handler for the Skribble backend.
// It upgrades HTTP connections to WebSocket connections, registers clients to rooms, and manages message broadcasting between clients in the same room.
package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nahom-zewdu/skribble-backend/internal/client"
	"github.com/nahom-zewdu/skribble-backend/internal/room"
	"github.com/nahom-zewdu/skribble-backend/pkg/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var roomManager = room.NewManager()

func (s *HTTPServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	name := r.URL.Query().Get("name")
	roomID := r.URL.Query().Get("room")

	if name == "" || roomID == "" {
		conn.Close()
		return
	}

	clientID := utils.GenerateID()
	c := client.NewClient(clientID, name, conn)
	c.RoomID = roomID

	room := roomManager.GetOrCreateRoom(roomID)
	room.Register(c)

	go c.WritePump()

	c.ReadPump(
		func(message []byte) {
			room.Broadcast(message)
		},
		func() {
			room.Unregister(c)
		},
	)
}
