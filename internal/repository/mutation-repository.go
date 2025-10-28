package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/IgorGrieder/Leaky-Bucket/internal/database"
)

type MutationRepository struct {
	Postgress *sql.DB
}

func NewMutationRepository(pg *sql.DB) *MutationRepository {
	return &MutationRepository{Postgress: pg}
}

func (repository *MutationRepository) QueryPixKey(mutation string, ctx context.Context) ([]database.MutationEntity, error) {
	ctxDb, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	const query = "SELECT pix_key FROM USERS WHERE $1"

	var entities []database.MutationEntity
	rows, err := repository.Postgress.QueryContext(ctxDb, query, mutation)
	if err != nil {
		return nil, fmt.Errorf("error executing the query %s in the database: %v", query, err)
	}

	for rows.Next() {
		var mutation database.MutationEntity

		if err := rows.Scan(&mutation.Key); err != nil {
			return nil, fmt.Errorf("error while scaning postgres row: %v", err)
		}

		entities = append(entities, mutation)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over query result: %v", err)
	}

	if len(entities) == 0 {
		return nil, fmt.Errorf("no PIX key related to the user: %w", sql.ErrNoRows)
	}

	return entities, nil
}
