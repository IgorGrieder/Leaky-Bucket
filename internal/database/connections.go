package database

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

func StartConns() (*redis.Client, *sql.DB) {
	return SetupRedis(), SetupPG()
}
