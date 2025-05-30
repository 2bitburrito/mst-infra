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
	fmt.Printf("Port: %s\n", cfg.Port)
	fmt.Printf("ApiKey: %s\n", cfg.ApiKey)
	fmt.Printf("DB URL: %s\n", cfg.DB.URL)
	fmt.Printf("CognitoPoolID: %s\n", cfg.CognitoPoolID)
	fmt.Printf("ReaperDuration: %s\n", cfg.ReaperDuration)

	return cfg, nil
}
