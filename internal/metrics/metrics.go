// internal/metrics/metrics.go
// This file defines the Metrics struct and related functions for tracking various metrics in the Skribble backend, such as active connections, rooms, messages, and bytes transferred.
// It uses atomic operations to ensure thread safety when updating these metrics.
package metrics

import (
	"sync/atomic"
)

// Metrics holds various metrics for the Skribble backend, including active connections, rooms, messages, and bytes transferred. It uses atomic operations to ensure thread safety when updating these metrics.
type Metrics struct {
	ActiveConnections int64 `json:"activeConnections"`
	PeakConnections   int64 `json:"peakConnections"`

	ActiveRooms int64 `json:"activeRooms"`
	PeakRooms   int64 `json:"peakRooms"`

	TotalMessages int64 `json:"totalMessages"`
	TotalBytes    int64 `json:"totalBytes"`

	DrawMessages int64 `json:"drawMessages"`
	ChatMessages int64 `json:"chatMessages"`
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

// IncRooms increments the count of active rooms and updates the peak rooms if the current count exceeds the previous peak.
func IncRooms() {
	current := atomic.AddInt64(
		&M.ActiveRooms,
		1,
	)

	updatePeak(
		&M.PeakRooms,
		current,
	)
}

func DecRooms() {
	atomic.AddInt64(
		&M.ActiveRooms,
		-1,
	)
}

// AddMessage increments the total message count, total bytes, and optionally the draw and chat message counts based on the provided parameters.
// It uses atomic operations to ensure thread safety when updating these metrics.
func AddMessage(size int, isDraw bool, isChat bool) {
	atomic.AddInt64(
		&M.TotalMessages,
		1,
	)

	atomic.AddInt64(
		&M.TotalBytes,
		int64(size),
	)

	if isDraw {
		atomic.AddInt64(
			&M.DrawMessages,
			1,
		)
	}

	if isChat {
		atomic.AddInt64(
			&M.ChatMessages,
			1,
		)
	}
}

// updatePeak updates the peak value if the current value exceeds the previous peak.
// It uses atomic operations to ensure thread safety when updating the peak value.
func updatePeak(
	peak *int64,
	current int64,
) {
	for {
		old := atomic.LoadInt64(peak)

		if current <= old {
			return
		}

		if atomic.CompareAndSwapInt64(
			peak,
			old,
			current,
		) {
			return
		}
	}
}
