package repository

import "database/sql"

type MutationRepository struct {
	Postgress *sql.DB
}

func NewMutationRepository(pg *sql.DB) *MutationRepository {
	return &MutationRepository{Postgress: pg}
}
