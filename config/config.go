package config

import (
	"context"
	"log"
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
	URL string
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
	dbUrl, err := getDbUrl()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Port:   os.Getenv("PORT"),
		ApiKey: os.Getenv("API_KEY"),
		DB: PostgresConfig{
			URL: dbUrl,
		},
		EmailClient:    sesClient,
		CognitoPoolID:  os.Getenv("COGNITO_POOL_ID"),
		ReaperDuration: 10 * time.Minute,
	}
	log.Println("Creating config...")
	log.Println("Environment running as:", os.Getenv("ENV"))
	log.Printf("Port: %v\n", len(cfg.Port) != 0)
	log.Printf("ApiKey: %v\n", len(cfg.ApiKey) != 0)
	log.Printf("DB URL: %v\n", len(cfg.DB.URL) != 0)
	log.Printf("ReaperDuration: %s\n", cfg.ReaperDuration)
	log.Printf("CognitoPoolID: %v\n", len(cfg.CognitoPoolID) != 0)

	return cfg, nil
}
