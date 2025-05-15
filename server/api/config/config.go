package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	DB             PostgresConfig
	Port           string
	ApiKey         string
	CognitoPoolID  string
	ReaperDuration time.Duration
}

type PostgresConfig struct {
	Username string
	Password string
	URL      string
	Port     string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't load env")
	}
	cfg := &Config{
		Port:   os.Getenv("PORT"),
		ApiKey: os.Getenv("API_KEY"),
		DB: PostgresConfig{
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PWD"),
			URL:      os.Getenv("DB_URL"),
			Port:     os.Getenv("DB_PORT"),
		},
		CognitoPoolID: os.Getenv("COGNITO_POOL_ID"),
	}
	return cfg, nil
}
