package database

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func SetupRedis() *redis.Client {
	const REDIS_ADDR = "redis"
	const REDIS_PORT = 6379

	redis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", REDIS_ADDR, REDIS_PORT),
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	err := redis.Ping(context.Background()).Err()
	if err != nil {
		fmt.Printf("Ending the execution %v", err)
		os.Exit(1)
	}

	return redis
}
