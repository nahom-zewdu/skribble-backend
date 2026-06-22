// internal/metrics/persistence.go
// This file defines the functions for persisting metrics to Redis.
// It includes an initialization function for setting up the Redis client and a flush loop that periodically saves the current metrics to Redis.
package metrics

import (
	"context"
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
