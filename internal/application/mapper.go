package application

import (
	"github.com/IgorGrieder/Leaky-Bucket/internal/database"
	"github.com/IgorGrieder/Leaky-Bucket/internal/domain"
)

func ToMutationAPI(entity database.MutationEntity) domain.Mutation {
	return domain.Mutation{
		PIX_KEY: entity.Key,
	}
}

func ToMutationAPISlice(entities []database.MutationEntity) []domain.Mutation {
	response := make([]domain.Mutation, 0, len(entities))
	for _, entity := range entities {
		response = append(response, ToMutationAPI(entity))
	}
	return response
}
