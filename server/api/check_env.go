package main

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func CheckEnv() (pqConnectionSting string, error error) {
	godotenv.Load()
	DB_URL := os.Getenv("DB_URL")
	if len(DB_URL) == 0 {
		return "", errors.New("DB_URL is not set")
	}

	return DB_URL, nil
}
