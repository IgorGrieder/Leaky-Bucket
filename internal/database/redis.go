package database

import (
	"context"
	"fmt"
	"os"

	"github.com/IgorGrieder/Leaky-Bucket/internal/config"
	"github.com/redis/go-redis/v9"
)

func SetupRedis(cfg *config.Config) *redis.Client {

	redis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.REDIS_ADDR, cfg.REDIS_PORT),
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	err := redis.Ping(context.Background()).Err()
	if err != nil {
		fmt.Printf("ending the execution %v", err)
		os.Exit(1)
	}

	return redis
}
