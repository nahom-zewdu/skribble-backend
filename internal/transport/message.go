// internal/transport/message.go
// This file defines the Message struct, which represents a message that can be sent between the server and clients in the Skribble backend.
// It includes a Type field to indicate the type of message and a Data field to hold the message payload.
package transport

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
