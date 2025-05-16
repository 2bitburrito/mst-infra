package main

import (
	"api/config"
	"api/jwt"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (api *API) createLoginCode(w http.ResponseWriter, r *http.Request) {
	cfg, _ := config.LoadConfig()
	var user User

	fmt.Println("Creating Login Code")

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
		"otc": otc,
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
