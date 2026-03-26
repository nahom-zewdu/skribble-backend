// internal/transport/message.go
// This file defines the Message struct, which represents a message that can be sent between the server and clients in the Skribble backend.
// It includes a Type field to indicate the type of message and a Data field to hold the message payload.
package transport

import "encoding/json"

/*
Message is the base envelope for all websocket communication.

Every message MUST follow this structure.
*/
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

/*
Incoming client message structure
*/
type ClientMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

/*
Chat message payload
*/
type ChatMessage struct {
	Text string `json:"text"`
}

/*
System message payload
*/
type SystemMessage struct {
	Text string `json:"text"`
}

/*
Drawing message payloads
*/
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

/*
Drawing actions
*/
type DrawStart struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type DrawMove struct {
	Points []Point `json:"points"`
}

/*
Drawing end and canvas clear actions
*/
type DrawEnd struct{}
type ClearCanvas struct{}
