package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"laboratory-internet-ai-test/config"
	"log"
	"time"
)

func Open(cfg config.Redis) (*redis.Client, error) {

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.Pass,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   2,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Printf("Redis connected successfully to %s ", addr)
	return client, nil
}
