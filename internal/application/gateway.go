package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IgorGrieder/Leaky-Bucket/internal/domain"
	"github.com/IgorGrieder/Leaky-Bucket/internal/repository"
)

type ProcessorService struct {
	LimitingRepository *repository.LimitingRepository
	MutationRepository *repository.MutationRepository
}

func (p *ProcessorService) ProcessMutation(mutation domain.Mutation, ctx context.Context) ([]domain.Mutation, error) {
	consumed, err := p.LimitingRepository.TryConsumeToken(ctx, "hi")

	if !consumed {
		return nil, &NoTokensError{}
	}

	entities, err := p.MutationRepository.QueryPixKey(mutation.PIX_KEY, ctx)

	if err != nil {
		log.Printf("error fetching pix keys check: %v", err)

		if errors.Is(err, sql.ErrNoRows) {
			return []domain.Mutation{}, nil
		}

		return nil, err
	}

	p.LimitingRepository.RefillToken(ctx, "hi")

	return ToMutationAPISlice(entities), nil
}

func (p *ProcessorService) FetchAndRefilTokens(ctx context.Context) error {
	var cursor uint64
	var err error

	for {
		var keys []string
		ctxRedis, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		keys, cursor, err = p.LimitingRepository.Redis.Scan(ctxRedis, cursor, "*", 50).Result()
		if err != nil {
			log.Printf("failed to SCAN keys from Redis: %v", err)
			return err
		}

		for _, key := range keys {
			err := p.refillToken(ctx, key)
			if err != nil {
				log.Printf("Failed to refill token for key '%s': %v", key, err)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (p *ProcessorService) refillToken(ctx context.Context, key string) error {
	err := p.LimitingRepository.RefillToken(ctx, key)
	if err != nil {
		return fmt.Errorf("error while refilling token %v", err)
	}

	return nil
}
