package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type LimitingRepository struct {
	Redis        *redis.Client
	MAX_ATTEMPTS int32

	locks   map[string]*sync.Mutex
	locksMu sync.Mutex
}

func NewLimitingRepository(redis *redis.Client) *LimitingRepository {
	return &LimitingRepository{
		Redis:        redis,
		MAX_ATTEMPTS: 10,
		locks:        make(map[string]*sync.Mutex),
	}
}

func (r *LimitingRepository) getLock(key string) *sync.Mutex {
	r.locksMu.Lock()
	mu, ok := r.locks[key]
	if !ok {
		mu = &sync.Mutex{}
		r.locks[key] = mu
	}
	r.locksMu.Unlock()
	return mu
}

func (r *LimitingRepository) QueryToken(ctx context.Context, key string) (int32, error) {
	mu := r.getLock(key)
	mu.Lock()
	defer mu.Unlock()

	ctxRedis, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	val, err := r.Redis.Get(ctxRedis, key).Int64()
	if err == nil {
		return int32(val), nil
	}
	if !errors.Is(err, redis.Nil) {
		return 0, fmt.Errorf("Error fetching tokens in the bucket: %v", err)
	}

	ctxSet, cancelSet := context.WithTimeout(ctx, 1*time.Second)
	defer cancelSet()

	attemptsLeft := r.MAX_ATTEMPTS
	wasSet, err := r.Redis.SetNX(ctxSet, key, attemptsLeft, 24*time.Hour).Result()
	if err != nil {
		return 0, fmt.Errorf("Error creating token key in the bucket: %v", err)
	}
	if wasSet {
		return attemptsLeft, nil
	}

	// Race lost — read again
	ctxGetAgain, cancelGetAgain := context.WithTimeout(ctx, 1*time.Second)
	defer cancelGetAgain()

	val, err = r.Redis.Get(ctxGetAgain, key).Int64()
	if err != nil {
		return 0, fmt.Errorf("Error fetching token after failed SetNX: %v", err)
	}
	return int32(val), nil
}

func (r *LimitingRepository) DecrementToken(ctx context.Context, key string) (int64, error) {
	mu := r.getLock(key)
	mu.Lock()
	defer mu.Unlock()

	ctxRedis, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	val, err := r.Redis.Decr(ctxRedis, key).Result()
	if err != nil {
		return 0, err
	}
	return val, nil
}
