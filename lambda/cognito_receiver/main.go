package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
)

type CognitoEvent struct {
	Version       string            `json:"version"`
	TriggerSource string            `json:"triggerSource"`
	Region        string            `json:"region"`
	UserPoolID    string            `json:"userPoolId"`
	UserName      string            `json:"userName"`
	CallerContext map[string]string `json:"callerContext"`
	Request       struct {
		UserAttributes UserAttributes `json:"userAttributes"`
	} `json:"request"`
	Response struct{} `json:"response"`
}

type UserAttributes struct {
	Sub                string `json:"sub"`
	Email              string `json:"email"`
	Name               string `json:"name"`
	ConfirmationStatus string `json:"user_status"`
}
type UserFromDB struct {
	Id                 string
	Email              string
	HasLicencse        bool
	EmailSubscription  bool
	CreatedAt          string
	NumberOfLicenses   string
	SubscribedToEmails bool
}

func createDbString() (pqConnectionSting string, error error) {
	DB_URL_WRITE := os.Getenv("DB_URL_WRITE")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_NAME := os.Getenv("DB_NAME")
	if len(DB_URL_WRITE) == 0 {
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
	pqConnectionSting = fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable", DB_URL_WRITE, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	return pqConnectionSting, nil
}

func handler(ctx context.Context, event CognitoEvent) (CognitoEvent, error) {
	fmt.Printf("received event: %v\n", event)
	fmt.Println("User Attributes:", event.Request.UserAttributes)

	user := event.Request.UserAttributes
	fmt.Println("NAME:", user.Name)

	if event.TriggerSource != "PostConfirmation_ConfirmSignUp" {
		fmt.Println("Trigger source is not PostConfirmation_ConfirmSignUp")
		return event, nil
	}
	err := db.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
	} else {
		fmt.Println("Database ping successful")
	}
	statement := `
    INSERT INTO users (id, email, full_name, has_license, email_subscription) 
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id;`

	var id string
	err = db.QueryRow(statement, event.UserName, user.Email, user.Name, false, false).Scan(&id)
	if err != nil {
		fmt.Println("Error in exec statement:", err)
		return event, err
	}
	fmt.Println("User inserted with ID:", id)

	fmt.Println("End of function")
	return CognitoEvent{}, nil
}

var db *sql.DB

func main() {
	fmt.Println("main called")
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

	defer db.Close()
	fmt.Println("Database Connection established")

	lambda.Start(handler)
}
