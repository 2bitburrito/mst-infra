package main

import (
	"api/config"
	"api/jwt"
	"encoding/json"
	"fmt"
	"net/http"
)

func createLoginCode(w http.ResponseWriter, r *http.Request) {
	cfg, _ := config.LoadConfig()
	var user User

	fmt.Println("Creating Login Code")

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	// Verify the jwt with cognito
	jwt.Verify(cfg.CognitoPoolID, user.Id, user.JWT)
}
