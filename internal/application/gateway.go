package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

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

func (p *ProcessorService) FetchAndRefilTokens() error {
	ctx := context.Background()

	var cursor uint64
	keys, cursor, err := p.LimitingRepository.Redis.Scan(ctx, cursor, "*", 10).Result()
	err := p.refillTokens(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProcessorService) refillTokens(ctx context.Context) error {
	err := p.LimitingRepository.RefillToken(ctx, "user:*")
	if err != nil {
		return fmt.Errorf("error while refilling token %v", err)
	}

	return nil
}
