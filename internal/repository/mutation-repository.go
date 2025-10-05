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

func (repository *MutationRepository) QueryPixKey(mutation domain.Mutation, ctx context.Context) {
	ctxDb, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	rows := repository.Postgress.QueryRowContext(ctxDb, "SELECT * FROM USERS")

}
