package database

import (
	"database/sql"

	"github.com/IgorGrieder/Leaky-Bucket/internal/config"
	"github.com/redis/go-redis/v9"
)

type Connections struct {
	Redis *redis.Client
	PG    *sql.DB
}

func StartConns(cfg *config.Config) *Connections {
	return &Connections{
		Redis: SetupRedis(cfg),
		PG:    SetupPG(cfg),
	}
}
