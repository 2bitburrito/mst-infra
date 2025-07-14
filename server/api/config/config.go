package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/2bitburrito/mst-infra/email"
	"github.com/aws/aws-sdk-go-v2/config"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	DB             PostgresConfig
	Port           string
	ApiKey         string
	CognitoPoolID  string
	EmailClient    email.EmailSender
	ReaperDuration time.Duration
	context        context.Context
}

type PostgresConfig struct {
	Username string
	Password string
	URL      string
	Port     string
}

func LoadConfig() (*Config, error) {
	context := context.Background()
	awsCfg, err := config.LoadDefaultConfig(context, config.WithRegion("us-west-1"))
	if err != nil {
		return nil, err
	}
	sesClient := email.SesEmailClient{
		AwsCfg:       awsCfg,
		SendingEmail: "hello@metasoundtools.com",
	}

	cfg := &Config{
		Port:   os.Getenv("PORT"),
		ApiKey: os.Getenv("API_KEY"),
		DB: PostgresConfig{
			URL: os.Getenv("DB_URL"),
		},
		EmailClient:    sesClient,
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
