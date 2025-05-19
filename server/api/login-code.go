package main

import (
	"api/config"
	"api/jwt"
	"encoding/json"
	"log"
	"net/http"
)

type CreateLoginCodeRequest struct {
	UserID string `json:"id"`
	JWT    string `json:"jwt"`
}

type LoginCodeRequest struct {
	UserID       string `json:"userId"`
	OneTimeToken string `json:"otc"`
}

func (api *API) createLoginCode(w http.ResponseWriter, r *http.Request) {
	cfg, _ := config.LoadConfig()
	var user User

	log.Println("Creating Login Code")

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("error in Jwt Verify %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	// Verify the jwt with cognito
	verified, err := jwt.Verify(cfg.CognitoPoolID, user.Id, user.JWT)
	if err != nil || !verified {
		log.Printf("error in Jwt Verify %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	otc := api.verificationStore.New(user.Id)
	returnObj := map[string]string{
		"userId": user.Id,
		"otc":    otc,
	}

	returnData, err := json.Marshal(returnObj)
	if err != nil {
		log.Printf("error marshalling otc: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(returnData)
}

func (api *API) checkLoginCode(w http.ResponseWriter, r *http.Request) {
	log.Println("Checking Login Code")
	var request LoginCodeRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("error decoding json in checkLoginCode %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	token, err := api.verificationStore.Get(request.UserID)
	if err != nil {
		log.Printf("Error:  %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Got otc from store: %v", token)
	if token != request.OneTimeToken {
		log.Printf("Invalid Token. \tReceived: %v \tWant: %v", request.OneTimeToken, token)
		http.Error(w, "Invalid token received", http.StatusUnauthorized)
		return
	}
	log.Printf("Successful Otc Match")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
