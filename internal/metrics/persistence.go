// internal/metrics/persistence.go
// This file defines the functions for persisting metrics to Redis.
// It includes an initialization function for setting up the Redis client and a flush loop that periodically saves the current metrics to Redis.
package metrics

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var redisClient *redis.Client

func InitRedis(
	url string,
) {

	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Fatal(err)
	}

	redisClient = redis.NewClient(opt)

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("Redis connected")

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

// LoadPersisted loads the persisted metrics from Redis and updates the in-memory metrics accordingly.
// It retrieves each metric value from Redis using the Get method and updates the corresponding fields in the Metrics struct using atomic operations to ensure thread safety.
// If any of the keys are not found in Redis, it defaults to 0 for that metric.
// This function is typically called during application startup to restore the metrics state from the previous run.
func LoadPersisted() {

	if redisClient == nil {
		return
	}

	load := func(key string) int64 {
		val, err := redisClient.Get(
			ctx,
			key,
		).Int64()

		if err != nil {
			return 0
		}

		return val
	}

	atomic.StoreInt64(
		&M.PeakConnections,
		load("peak_connections"),
	)

	atomic.StoreInt64(
		&M.PeakRooms,
		load("peak_rooms"),
	)

	atomic.StoreInt64(
		&M.TotalMessages,
		load("total_messages"),
	)

	atomic.StoreInt64(
		&M.TotalBytes,
		load("total_bytes"),
	)

	atomic.StoreInt64(
		&M.DrawMessages,
		load("draw_messages"),
	)

	atomic.StoreInt64(
		&M.ChatMessages,
		load("chat_messages"),
	)

	log.Println("Loaded persisted metrics")
}
