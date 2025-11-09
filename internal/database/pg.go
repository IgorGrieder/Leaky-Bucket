package database

import (
	"fmt"
	"log/slog"
	"os"

	"database/sql"

	"github.com/IgorGrieder/Leaky-Bucket/internal/config"
	_ "github.com/lib/pq"
)

func SetupPG(cfg *config.Config) *sql.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.HOST, cfg.PORT_PG, cfg.USER, cfg.PG_PASS, cfg.DB_NAME)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		slog.Error("failed connecting into PG")
		os.Exit(1)
	}

	return db
}
