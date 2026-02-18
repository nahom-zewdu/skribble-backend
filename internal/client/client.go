// internal/client/client.go
// This file defines the Client struct, which represents a connected client in the Skribble backend.
// It includes fields for the client's ID, name, WebSocket connection, and a channel for sending messages to the client. The Client struct also has methods for reading messages from the client and writing messages to the client.

package client

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Name string

	Conn *websocket.Conn

	Send chan []byte // outbound messages

	RoomID string
}

func NewClient(id, name string, conn *websocket.Conn) *Client {
	return &Client{
		ID:   id,
		Name: name,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
}

func (c *Client) ReadPump(onMessage func([]byte), onClose func()) {
	defer func() {
		onClose()
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		onMessage(message)
	}
}

// WritePump listens on the Send channel and writes messages to the WebSocket connection.
// It should be run in a separate goroutine for each client.
// If an error occurs while writing, it logs the error and breaks the loop, which will eventually lead to closing the connection.
func (c *Client) WritePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
}
