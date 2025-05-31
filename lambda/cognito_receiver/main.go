package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

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

	if event.TriggerSource != "PostConfirmation_ConfirmSignUp" {
		log.Println("Trigger source is not PostConfirmation_ConfirmSignUp")
		return event, nil
	}

	apiKey := os.Getenv("API_KEY")
	if len(apiKey) == 0 {
		log.Println("apiKey not found")
		return event, errors.New("DB_URL is not set")
	}

	body, err := json.Marshal(user)
	if err != nil {
		log.Println("couldn't marshal user body")
		return event, err
	}

	request, err := http.NewRequest("POST", "https://meta-sound-tools.fly.dev/api/cognito-user", bytes.NewReader(body))
	if err != nil {
		log.Println("error creating new request")
		return event, err
	}

	request.Header.Set("X-API-Key", apiKey)
	request.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		log.Println("error sending HTTP Req")
		return event, err
	}
	defer response.Body.Close()

	log.Println("Response Status: ", response.Status)

	log.Println("End of function")
	return event, nil
}

func main() {
	log.Println("main called")

	lambda.Start(handler)
}
