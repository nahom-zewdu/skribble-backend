// utils/id.go
// This file defines a utility function to generate unique IDs for players and game sessions in the Skribble backend.

package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
