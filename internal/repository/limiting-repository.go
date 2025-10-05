package repository

import "github.com/redis/go-redis/v9"

type LimitingRepository struct {
	Redis *redis.Client
}

func NewLimitingRepository(redis *redis.Client) *LimitingRepository {
	return &LimitingRepository{Redis: redis}
}
