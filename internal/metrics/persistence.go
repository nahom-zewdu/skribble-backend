// internal/metrics/persistence.go
// This file defines the functions for persisting metrics to Redis.
// It includes an initialization function for setting up the Redis client and a flush loop that periodically saves the current metrics to Redis.
package metrics

import (
	"context"

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
