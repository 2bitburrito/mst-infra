package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

type CognitoEvent struct {
	Version       string         `json:"version"`
	TriggerSource string         `json:"triggerSource"`
	Region        string         `json:"region"`
	UserPoolID    string         `json:"userPoolId"`
	CallerContext map[string]any `json:"callerContext"`

	Request struct {
		UserAttributes UserAttributes `json:"userAttributes"`
		NewDeviceUsed  bool           `json:"newDeviceUsed"`
	} `json:"request"`
	ClientMetadata map[string]string `json:"clientMetadata"`
	Response       map[string]any    `json:"response"`
}

// User Attributes: map[cognito:user_status:CONFIRMED email:hughpalmerproduction@gmail.com email_verified:true sub:39b9692e-3061-70ff-29db-29125abe9c95]

type UserAttributes struct {
	Id                 string `json:"sub"`
	Email              string `json:"email"`
	ConfirmationStatus string `json:"user_status"`
}

func createDbString() (pqConnectionSting string, error error) {
	DB_URL := os.Getenv("DB_URL")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")
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

func handler(ctx context.Context, event CognitoEvent) CognitoEvent {
	fmt.Println("handler called")
	fmt.Println("Context:", ctx)
	fmt.Printf("received event: %v\n", event)
	fmt.Println("User Attributes:", event.Request.UserAttributes)

	err := db.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
	}
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println("Error in select query:", err)
	}
	fmt.Println("User rows:", rows)

	return event
}

var db *sql.DB

func main() {
	connectionString, err := createDbString()
	if err != nil {
		fmt.Println("Error creating connection string:", err)
	}
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("Error opening database:", err)
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(3)
	db.SetMaxOpenConns(3)

	fmt.Println("Database Connection established")

	lambda.Start(handler)
}
