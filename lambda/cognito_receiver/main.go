package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

type CognitoEvent struct {
	Version       string                 `json:"version"`
	TriggerSource string                 `json:"triggerSource"`
	Region        string                 `json:"region"`
	UserPoolID    string                 `json:"userPoolId"`
	CallerContext map[string]interface{} `json:"callerContext"`
	Request       struct {
		UserAttributes map[string]string `json:"userAttributes"`
		NewDeviceUsed  bool              `json:"newDeviceUsed"`
	} `json:"request"`
	Response map[string]interface{} `json:"response"`
}

func createDbString() (pqConnectionSting string, error error) {
	DB_URL := os.Getenv("DB_URL")
	DB_PORT := os.Getenv("DB_PORT")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_USER := os.Getenv("DB_USER")
	DB_NAME := os.Getenv("DB_NAME")
	if len(DB_URL) == 0 {
		return "", errors.New("DB_URL is not set")
	}
	if len(DB_PASSWORD) == 0 {
		return "", errors.New("DB_PASSWORD is not set")
	}
	if len(DB_USER) == 0 {
		return "", errors.New("DB_USER is not set")
	}
	if len(DB_NAME) == 0 {
		return "", errors.New("DB_NAME is not set")
	}

	pqConnectionSting = fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable", DB_URL, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	return pqConnectionSting, nil
}

func handler(ctx context.Context, event CognitoEvent) (CognitoEvent, error) {
	fmt.Printf("received event: %+v\n", event)
	return CognitoEvent{}, nil
}

func main() {
	lambda.Start(handler)
}
