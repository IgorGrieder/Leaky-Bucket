package database

import (
	"database/sql"

	"github.com/IgorGrieder/Leaky-Bucket/internal/config"
	"github.com/redis/go-redis/v9"
)

func StartConns(cfg *config.Config) (*redis.Client, *sql.DB) {
	return SetupRedis(cfg), SetupPG(cfg)
}
