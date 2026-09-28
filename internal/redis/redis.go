package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedis() (*redis.Client, error) {

	addr := os.Getenv("REDIS_ADDR")

	if addr == "" {
		addr = "localhost:6379"
	}


	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return client, nil
}