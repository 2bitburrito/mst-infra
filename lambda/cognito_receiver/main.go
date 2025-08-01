package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	_ "github.com/joho/godotenv"
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

type API struct {
	Key string
	URL string
}

type ServerErrror struct {
	Error string
}

func handleRequest(ctx context.Context, event CognitoEvent) (CognitoEvent, error) {
	api := &API{
		Key: os.Getenv("API_KEY"),
		URL: os.Getenv("API_URL"),
	}

	switch event.TriggerSource {
	case "PostConfirmation_ConfirmSignUp":
		return confirmSignup(ctx, event, api)
	default:
		return CognitoEvent{}, nil
	}
}

func confirmSignup(ctx context.Context, event CognitoEvent, api *API) (CognitoEvent, error) {
	url := fmt.Sprintf("%s/api/cognito-user", api.URL)
	user := event.Request.UserAttributes
	body, err := json.Marshal(user)
	if err != nil {
		return event, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return event, err
	}
	req.Header.Add("X-API-Key", api.Key)
	req.Header.Add("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return event, err
	}
	if resp.StatusCode != http.StatusOK {
		var serverErr ServerErrror
		json.NewDecoder(resp.Body).Decode(&serverErr)
		defer resp.Body.Close()
		return CognitoEvent{}, fmt.Errorf("%s", serverErr.Error)
	}
	return CognitoEvent{}, err
}

func main() {
	lambda.Start(handleRequest)
}
