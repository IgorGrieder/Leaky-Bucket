package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Define the Lua script.
// This script is the atomic "Get-or-Create" operation
var getTokenOrSetScript = redis.NewScript(`
-- KEYS[1] = key
-- ARGV[1] = default_value
-- ARGV[2] = ttl_seconds

-- Atomically set the key if it does not exist (NX), with an expiration (EX)
local was_set = redis.call('SET', KEYS[1], ARGV[1], 'EX', ARGV[2], 'NX')

if was_set then
    -- We successfully set it. Return the default value we just set.
    return ARGV[1]
else
    -- It already existed. Get and return its current value.
    return redis.call('GET', KEYS[1])
end
`)

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

	keys := []string{key}
	args := []any{
		r.MAX_ATTEMPTS,
		int((24 * time.Hour).Seconds()),
	}

	val, err := getTokenOrSetScript.Run(ctxRedis, r.Redis, keys, args...).Int64()
	if err != nil {
		// Handle redis.Nil in case the key expired between the SET fail and the GET
		if errors.Is(err, redis.Nil) {
			// If it's nil, it means it expired right as we checked.
			// We can just return the default, as the next call will set it.
			// Or, we could retry, but returning default is simpler.
			return r.MAX_ATTEMPTS, nil
		}
		return 0, fmt.Errorf("Error running GetOrSet script: %v", err)
	}

	return int32(val), nil
}

func (r *LimitingRepository) DecrementToken(ctx context.Context, key string) (int64, error) {
	ctxRedis, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	val, err := r.Redis.Decr(ctxRedis, key).Result()
	if err != nil {
		return 0, err
	}

	return val, nil
}
