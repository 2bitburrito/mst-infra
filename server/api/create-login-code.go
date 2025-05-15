package main

import (
	"api/config"
	"fmt"
	"net/http"

	cognitoJwtVerify "github.com/jhosan7/cognito-jwt-verify"
)

func createLoginCode(w http.ResponseWriter, r *http.Request) {
	cfg, _ := config.LoadConfig()

	fmt.Println("Creating Login Code")
	id := r.PathValue("id")
	jwt := r.PathValue("Authorization")
	// Verify the jwt with cognito
	fmt.Printf("User ID: %v", id)
	fmt.Printf("User JWT: %v", jwt)

	cognitoCfg := cognitoJwtVerify.Config{
		UserPoolId: cfg.CognitoPoolID,
		ClientId:   id,
	}

	verify, err := cognitoJwtVerify.Create(cognitoCfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	payload, err := verify.Verify("eyJraW...")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println(payload)
}
