package application

import (
	"context"
	"database/sql"
	"errors"

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
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.Mutation{}, nil
		}

		return nil, err
	}

	var response []domain.Mutation
	for _, entity := range entities {
		mappedMutation := domain.Mutation{
			PIX_KEY: entity.Key,
		}
		response = append(response, mappedMutation)
	}

	return response, nil
}
