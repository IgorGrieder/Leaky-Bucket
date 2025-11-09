package config

import (
	"os"
	"strconv"
)

type Config struct {
	PORT       int
	REDIS_ADDR string
	REDIS_PORT int
	HOST       string
	PORT_PG    int
	USER       string
	DB_NAME    string
	PG_PASS    string
	HASH       string
}

func NewConfig() *Config {
	port := parseInt(getEnv("PORT", "8080"))
	reddisAddr := getEnv("REDIS_ADDR", "localhost")
	reddisPort := parseInt(getEnv("REDIS_PORT", "6379"))
	host := getEnv("PG_HOST", "localhost")
	portPG := parseInt(getEnv("PG_PORT", "5432"))
	user := getEnv("PG_USER", "postgres")
	dbname := getEnv("PG_DB", "leaky_bucket")
	pgPass := getEnv("PG_PASS", "none")
	hash := getEnv("HASH", "733e6f1d9d33bd2e754da2ba56aeb563718fc68f5c139f527e96a13a1b2671fb")

	return &Config{
		PORT:       port,
		REDIS_ADDR: reddisAddr,
		REDIS_PORT: reddisPort,
		HOST:       host,
		PORT_PG:    portPG,
		USER:       user,
		DB_NAME:    dbname,
		PG_PASS:    pgPass,
		HASH:       hash,
	}
}

func getEnv(key string, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
