package database

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func SetupRedis() *redis.Client {
	const REDIS_ADDR = "redis"
	const REDIS_PORT = 6379

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", REDIS_ADDR, REDIS_PORT),
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	return redisClient
}
