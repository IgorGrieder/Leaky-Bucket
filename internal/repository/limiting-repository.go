package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type LimitingRepository struct {
	Redis        *redis.Client
	MAX_ATTEMPTS int32
}

func NewLimitingRepository(redis *redis.Client) *LimitingRepository {
	return &LimitingRepository{Redis: redis, MAX_ATTEMPTS: 10}
}

func (r *LimitingRepository) QueryToken(ctx context.Context, key string) (int32, error) {
	ctxRedis, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	val, err := r.Redis.Get(ctxRedis, key).Int64()

	if err != nil {

		if errors.Is(err, redis.Nil) {
			newValue, err := r.CreateToken(ctx, key)
			return newValue, err
		}

		return 0, fmt.Errorf("Error fetching tokens in the bucket %v", err)
	}

	return int32(val), nil
}

func (r *LimitingRepository) CreateToken(ctx context.Context, key string) (int32, error) {
	ctxRedis, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	attemptsLeft := r.MAX_ATTEMPTS

	err := r.Redis.Set(ctxRedis, key, attemptsLeft, 24*time.Hour).Err()
	if err != nil {
		return 0, fmt.Errorf("Error creating token key in the bucket %v", err)
	}

	return attemptsLeft, nil
}

func (r *LimitingRepository) DecrementToken(ctx context.Context, key string) error {
	ctxRedis, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	_, err := r.Redis.Decr(ctxRedis, key).Result()
	if err != nil {
		return err
	}

	return nil
}
