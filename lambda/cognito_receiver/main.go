package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
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

func handler(ctx context.Context, event CognitoEvent) (CognitoEvent, error) {
	log.Printf("received event: %v\n", event)
	log.Println("User Attributes:", event.Request.UserAttributes)

	user := event.Request.UserAttributes
	apiKey := os.Getenv("API_KEY")
	if len(apiKey) == 0 {
		log.Println("apiKey not found")
		return event, errors.New("DB_URL is not set")
	}

	if event.TriggerSource != "PostConfirmation_ConfirmSignUp" {
		log.Println("Trigger source is not PostConfirmation_ConfirmSignUp")
		return event, nil
	}
	log.Println("1")
	body, err := json.Marshal(&user)
	if err != nil {
		log.Println("couldn't marshal user body")
		return event, err
	}
	log.Println("2")

	request, err := http.NewRequest("POST", "https://meta-sound-tools.fly.dev/api/cognito-user", bytes.NewBuffer(body))
	if err != nil {
		log.Println("error creating new request")
		return event, err
	}
	log.Println("3")

	request.Header.Set("X-API-Key", apiKey)
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{Timeout: 6 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		log.Println("error sending HTTP Req")
		return event, err
	}
	defer response.Body.Close()
	log.Println("4")

	log.Println("Response Status: ", response.Status)

	log.Println("End of function")
	return event, nil
}

func main() {
	log.Println("main called")

	lambda.Start(handler)
}
