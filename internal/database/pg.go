package database

import (
	"fmt"
	"os"

	"database/sql"

	_ "github.com/lib/pq"
)

func SetupPG() *sql.DB {

	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "none"
		dbname   = "leaky-bucket"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("Ending the execution %v", err)
		os.Exit(1)
	}

	return db
}
