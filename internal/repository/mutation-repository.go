package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/IgorGrieder/Leaky-Bucket/internal/domain"
)

type MutationRepository struct {
	Postgress *sql.DB
}

func NewMutationRepository(pg *sql.DB) *MutationRepository {
	return &MutationRepository{Postgress: pg}
}

func (repository *MutationRepository) QueryPixKey(mutation string, ctx context.Context) ([]domain.MutationEntity, error) {
	ctxDb, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	var entities []domain.MutationEntity
	rows, err := repository.Postgress.QueryContext(ctxDb, "SELECT pix_key FROM USERS WHERE $1", mutation)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var mutation domain.MutationEntity

		if err := rows.Scan(&mutation.Key); err != nil {
			return nil, fmt.Errorf("Error while scaning postgres row: %v", err)
		}

		entities = append(entities, mutation)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error while iterating over query result: %v", err)
	}

	if len(entities) == 0 {
		return nil, fmt.Errorf("No PIX key related to the user: %w", sql.ErrNoRows)
	}

	return entities, nil
}
