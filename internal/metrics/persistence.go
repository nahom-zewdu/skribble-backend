// internal/metrics/persistence.go
// This file defines the functions for persisting metrics to Redis.
// It includes an initialization function for setting up the Redis client and a flush loop that periodically saves the current metrics to Redis.
package metrics

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var redisClient *redis.Client

func InitRedis(
	url string,
	password string,
) {

	redisClient = redis.NewClient(
		&redis.Options{
			Addr:     url,
			Password: password,
		},
	)

	go flushLoop()
}

// flushLoop runs an infinite loop that flushes the current metrics to Redis every 15 seconds. It uses a ticker to trigger the flush operation at regular intervals. If the Redis client is not initialized, the Flush function will simply return without doing anything.
func flushLoop() {
	ticker := time.NewTicker(
		15 * time.Second,
	)

	defer ticker.Stop()

	for range ticker.C {
		Flush()
	}
}

// Flush saves the current metrics to Redis. It uses a pipeline to set multiple keys in a single round trip for efficiency. If the Redis client is not initialized, it simply returns without doing anything.
func Flush() {

	if redisClient == nil {
		return
	}

	m := Snapshot()

	pipe := redisClient.Pipeline()

	pipe.Set(
		ctx,
		"peak_connections",
		m.PeakConnections,
		0,
	)

	pipe.Set(
		ctx,
		"peak_rooms",
		m.PeakRooms,
		0,
	)

	pipe.Set(
		ctx,
		"total_messages",
		m.TotalMessages,
		0,
	)

	pipe.Set(
		ctx,
		"total_bytes",
		m.TotalBytes,
		0,
	)

	pipe.Set(
		ctx,
		"draw_messages",
		m.DrawMessages,
		0,
	)

	pipe.Set(
		ctx,
		"chat_messages",
		m.ChatMessages,
		0,
	)

	_, _ = pipe.Exec(ctx)
}

// Snapshot creates a copy of the current metrics by loading the values atomically.
// This ensures that the snapshot is consistent and reflects the state of the metrics at the time it was taken.
// The returned Metrics struct contains the current values of all the metrics, which can then be used for persistence or reporting purposes.
func Snapshot() Metrics {
	return Metrics{
		ActiveConnections: atomic.LoadInt64(
			&M.ActiveConnections,
		),

		PeakConnections: atomic.LoadInt64(
			&M.PeakConnections,
		),

		ActiveRooms: atomic.LoadInt64(
			&M.ActiveRooms,
		),

		PeakRooms: atomic.LoadInt64(
			&M.PeakRooms,
		),

		TotalMessages: atomic.LoadInt64(
			&M.TotalMessages,
		),

		TotalBytes: atomic.LoadInt64(
			&M.TotalBytes,
		),

		DrawMessages: atomic.LoadInt64(
			&M.DrawMessages,
		),

		ChatMessages: atomic.LoadInt64(
			&M.ChatMessages,
		),
	}
}
