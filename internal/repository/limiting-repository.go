package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type LimitingRepository struct {
	Redis        *redis.Client
	MAX_ATTEMPTS int32
}

func NewLimitingRepository(redis *redis.Client) *LimitingRepository {
	return &LimitingRepository{
		Redis:        redis,
		MAX_ATTEMPTS: 10,
	}
}

func (r *LimitingRepository) TryConsumeToken(ctx context.Context, key string) (bool, error) {
	ctxSet, cancelSet := context.WithTimeout(ctx, 1*time.Second)
	defer cancelSet()

	wasSet, err := r.Redis.SetNX(ctxSet, key, r.MAX_ATTEMPTS-1, 24*time.Hour).Result()
	if err != nil {
		return false, fmt.Errorf("error initializing token bucket: %w", err)
	}

	if wasSet {
		return true, nil
	}

	ctxDecr, cancelDecr := context.WithTimeout(ctx, 1*time.Second)
	defer cancelDecr()

	remaining, err := r.Redis.Decr(ctxDecr, key).Result()
	if err != nil {
		return false, fmt.Errorf("error decrementing token bucket: %w", err)
	}

	if remaining < 0 {
		ctxRollback, cancelRollback := context.WithTimeout(ctx, 1*time.Second)
		defer cancelRollback()

		r.Redis.Incr(ctxRollback, key)

		return false, nil
	}

	return true, nil
}

func (r *LimitingRepository) RefillToken(ctx context.Context, key string) error {
	ctxRedis, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	newValue, err := r.Redis.Incr(ctxRedis, key).Result()
	if err != nil {
		return fmt.Errorf("error refilling token: %w", err)
	}

	// Ensure limit
	if newValue > int64(r.MAX_ATTEMPTS) {
		ctxSet, cancelSet := context.WithTimeout(ctx, 1*time.Second)
		defer cancelSet()

		r.Redis.Set(ctxSet, key, r.MAX_ATTEMPTS, 24*time.Hour)
	}

	return nil
}
