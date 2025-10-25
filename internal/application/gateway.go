package application

import (
	"context"

	"github.com/IgorGrieder/Leaky-Bucket/internal/domain"
	"github.com/IgorGrieder/Leaky-Bucket/internal/repository"
)

type ProcessorService struct {
	LimitingRepository *repository.LimitingRepository
	MutationRepository *repository.MutationRepository
}

func (p *ProcessorService) ProcessMutation(mutation domain.Mutation, ctx context.Context) (domain.Mutation, error) {
	return domain.Mutation{PIX_KEY: "Hello"}, nil
}
