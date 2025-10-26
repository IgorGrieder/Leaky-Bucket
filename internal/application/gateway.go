package application

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/IgorGrieder/Leaky-Bucket/internal/domain"
	"github.com/IgorGrieder/Leaky-Bucket/internal/repository"
)

type ProcessorService struct {
	LimitingRepository *repository.LimitingRepository
	MutationRepository *repository.MutationRepository
}

func (p *ProcessorService) ProcessMutation(mutation domain.Mutation, ctx context.Context) ([]domain.Mutation, error) {
	entities, err := p.MutationRepository.QueryPixKey(mutation.PIX_KEY, ctx)

	if err != nil {
		log.Printf("Error fetching pix keys check: %v", err)

		if errors.Is(err, sql.ErrNoRows) {
			return []domain.Mutation{}, nil
		}

		return nil, err
	}

	return ToMutationAPISlice(entities), nil
}
