package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/IgorGrieder/Leaky-Bucket/internal/domain"
)

type MutationRepository struct {
	Postgress *sql.DB
}

func NewMutationRepository(pg *sql.DB) *MutationRepository {
	return &MutationRepository{Postgress: pg}
}

func (repository *MutationRepository) QueryPixKey(mutation domain.Mutation, ctx context.Context) ([]domain.MutationEntity, error) {
	ctxDb, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	var entities []domain.MutationEntity
	rows, err := repository.Postgress.QueryContext(ctxDb, "SELECT * FROM USERS WHERE ?", mutation.PIX_KEY)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var mutation domain.MutationEntity

		if err := rows.Scan(&mutation.Key); err != nil {
			return nil, err
		}

		entities = append(entities, mutation)
	}

	return entities, nil

}
