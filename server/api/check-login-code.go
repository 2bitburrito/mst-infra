package main

import (
	"api/config"
	"encoding/json"
	"log"
	"net/http"
)

func (api *API) checkLoginCode(w http.ResponseWriter, r *http.Request) {
	cfg, _ := config.LoadConfig()
	log.Println("Checking Login Code")
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("error decoding json in checkLoginCode %v", err.Error())
	}
}
