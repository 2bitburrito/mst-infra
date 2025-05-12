package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type CognitoObject struct {
	ClientId     string `json:"ClientId"`
	ClientSecret string `json:"ClientSecret"`
	Token        string `json:"Token"`
}

// NOTE: THE curl cmd from AWS docs:
// curl --location 'auth.mydomain.com/oauth2/revoke' \
// --header 'Content-Type: application/x-www-form-urlencoded' \
// --header 'Authorization: Basic Base64Encode(client_id:client_secret)' \
// --data-urlencode 'token=abcdef123456789ghijklexample' \
// --data-urlencode 'client_id=1example23456789'

func RevokeAccess(request CognitoObject) error {
	endpoint := "https://metasoundtools.auth.us-west-1.amazoncognito.com/oauth2/revoke"

	auth := base64.StdEncoding.EncodeToString([]byte(request.ClientId + ":" + request.ClientSecret))

	formData := url.Values{}
	formData.Set("token", request.Token)
	formData.Set("client_id", request.ClientId)

	reqBody := strings.NewReader(formData.Encode())
	req, err := http.NewRequest("POST", endpoint, reqBody)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic"+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to revoke access in cognito call:%d - %v", resp.StatusCode, resp.Status)
	}

	fmt.Printf("response body: %v\n", resp.Body)
	return nil
}
