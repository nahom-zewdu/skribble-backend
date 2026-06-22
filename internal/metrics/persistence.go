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
