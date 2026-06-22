// internal/server/websocket.go
// This file defines the WebSocket handler for the Skribble backend.
// It upgrades HTTP connections to WebSocket connections, registers clients to rooms, and manages message broadcasting between clients in the same room.
package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nahom-zewdu/skribble-backend/internal/client"
	"github.com/nahom-zewdu/skribble-backend/internal/metrics"
	"github.com/nahom-zewdu/skribble-backend/internal/pkg/utils"
	"github.com/nahom-zewdu/skribble-backend/internal/room"
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
	mode := r.URL.Query().Get("mode")
	roomID := r.URL.Query().Get("room")

	if name == "" || mode == "" {
		conn.Close()
		return
	}

	var room *room.Room

	switch mode {

	case "public":
		room = roomManager.FindOrCreatePublicRoom()

	case "private_create":
		room = roomManager.CreatePrivateRoom()

	case "private_join":
		if roomID == "" {
			conn.Close()
			return
		}

		privateRoom, ok := roomManager.GetRoom(roomID)
		if !ok {
			conn.Close()
			return
		}

		if !privateRoom.IsJoinable() {
			conn.Close()
			return
		}

		room = privateRoom

	default:
		conn.Close()
		return
	}

	clientID := utils.GenerateID()

	c := client.NewClient(clientID, name, conn)
	c.RoomID = room.ID

	metrics.IncConnections() // Increment the active connections metric when a new client connects
	room.Register(c)

	go c.WritePump()

	c.ReadPump(
		func(message []byte) {
			room.HandleClientMessage(c, message)
		},
		func() {
			room.Unregister(c)
		},
	)
}
