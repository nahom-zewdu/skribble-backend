// internal/metrics/metrics.go
// This file defines the Metrics struct and related functions for tracking various metrics in the Skribble backend, such as active connections, rooms, messages, and bytes transferred.
// It uses atomic operations to ensure thread safety when updating these metrics.
package metrics

import (
	"sync/atomic"
)

// Metrics holds various metrics for the Skribble backend, including active connections, rooms, messages, and bytes transferred. It uses atomic operations to ensure thread safety when updating these metrics.
type Metrics struct {
	ActiveConnections int64
	PeakConnections   int64

	ActiveRooms int64
	PeakRooms   int64

	TotalMessages int64
	TotalBytes    int64

	DrawMessages int64
	ChatMessages int64
}

var M = &Metrics{}

// IncConnections increments the count of active connections and updates the peak connections if the current count exceeds the previous peak.
func IncConnections() {
	current := atomic.AddInt64(
		&M.ActiveConnections,
		1,
	)

	updatePeak(
		&M.PeakConnections,
		current,
	)
}

// DecConnections decrements the count of active connections.
func DecConnections() {
	atomic.AddInt64(
		&M.ActiveConnections,
		-1,
	)
}
