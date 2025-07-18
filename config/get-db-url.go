package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func getDbUrl() (pqConnectionSting string, error error) {
	godotenv.Load()
	ENV := os.Getenv("ENV")
	DEV_DB_URL := os.Getenv("DEV_DB_URL")
	DB_URL := os.Getenv("DB_URL")

	if ENV == "dev" {
		return DEV_DB_URL, nil
	}
	if len(DB_URL) == 0 {
		return "", errors.New("DB_URL is not set")
	}

	return DB_URL, nil
}
