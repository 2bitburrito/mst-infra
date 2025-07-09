package config

import (
	"fmt"
	"os"
	"time"

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
	cfg := &Config{
		Port:   os.Getenv("PORT"),
		ApiKey: os.Getenv("API_KEY"),
		DB: PostgresConfig{
			URL: os.Getenv("DB_URL"),
		},
		CognitoPoolID:  os.Getenv("COGNITO_POOL_ID"),
		ReaperDuration: 10 * time.Minute,
	}
	fmt.Printf("Port: %v\n", len(cfg.Port) != 0)
	fmt.Printf("ApiKey: %v\n", len(cfg.ApiKey) != 0)
	fmt.Printf("DB URL: %v\n", len(cfg.DB.URL) != 0)
	fmt.Printf("CognitoPoolID: %v\n", len(cfg.CognitoPoolID) != 0)
	fmt.Printf("ReaperDuration: %s\n", cfg.ReaperDuration)

	return cfg, nil
}
