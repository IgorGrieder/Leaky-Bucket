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
	return &LimitingRepository{
		Redis:        redis,
		MAX_ATTEMPTS: 10,
	}
}

func (r *LimitingRepository) TryConsumeToken(ctx context.Context, key string) (bool, int32, error) {
	ctxRedis, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	remaining, err := r.Redis.Decr(ctxRedis, key).Result()

	if errors.Is(err, redis.Nil) || err != nil {
		ctxInit, cancelInit := context.WithTimeout(ctx, 1*time.Second)
		defer cancelInit()

		// SetNX = só cria se não existir (evita race na criação)
		wasSet, err := r.Redis.SetNX(ctxInit, key, r.MAX_ATTEMPTS-1, 24*time.Hour).Result()
		if err != nil {
			return false, 0, fmt.Errorf("error initializing token bucket: %w", err)
		}

		if wasSet {
			return true, r.MAX_ATTEMPTS - 1, nil
		}

		// Outra goroutine criou primeiro - tenta de novo
		return r.TryConsumeToken(ctx, key)
	}

	// Ficou negativo = sem tokens disponíveis
	if remaining < 0 {
		// Rollback: devolve o token
		ctxRollback, cancelRollback := context.WithTimeout(ctx, 1*time.Second)
		defer cancelRollback()

		r.Redis.Incr(ctxRollback, key)
		return false, 0, nil
	}

	// Sucesso - consumiu 1 token
	return true, int32(remaining), nil
}

func (r *LimitingRepository) RefillToken(ctx context.Context, key string) error {
	ctxRedis, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	newValue, err := r.Redis.Incr(ctxRedis, key).Result()
	if err != nil {
		return fmt.Errorf("error refilling token: %w", err)
	}

	// Se ultrapassou o limite, ajusta para MAX
	if newValue > int64(r.MAX_ATTEMPTS) {
		ctxSet, cancelSet := context.WithTimeout(ctx, 1*time.Second)
		defer cancelSet()

		r.Redis.Set(ctxSet, key, r.MAX_ATTEMPTS, 24*time.Hour)
	}

	return nil
}

