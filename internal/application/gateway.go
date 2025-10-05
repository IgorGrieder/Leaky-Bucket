package application

import (
	"context"

	"github.com/IgorGrieder/Leaky-Bucket/internal/domain"
)

type ProcessorService struct {
}

func (p *ProcessorService) ProcessMutation(mutation domain.Mutation, ctx context.Context) error {
	return nil
}
