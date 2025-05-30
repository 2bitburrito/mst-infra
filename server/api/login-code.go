package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/2bitburrito/mst-infra/server/api/config"
	"github.com/2bitburrito/mst-infra/server/api/jwt"
	"github.com/2bitburrito/mst-infra/server/api/utils"
)

type CreateLoginCodeRequest struct {
	UserID string `json:"id"`
	JWT    string `json:"jwt"`
}

type LoginCodeRequest struct {
	UserID       *string `json:"userId,omitempty"`
	OneTimeToken string  `json:"otc"`
	MachineID    string  `json:"machineId"`
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
	verified, err := jwt.VerifyCognitoJWT(cfg.CognitoPoolID, user.Id, user.JWT)
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

	var token string

	if request.UserID == nil {
		// Get Token from store from OTC
		var userId string
		userId, token, err = api.verificationStore.GetFromOTC(request.OneTimeToken)
		if err != nil {
			log.Printf("Error:  %v", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		request.UserID = &userId
	} else {
		// Get Token from store matching userID:
		token, err = api.verificationStore.Get(*request.UserID)
		if err != nil {
			log.Printf("Error:  %v", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	log.Printf("Retrieved OTC from store: %v", token)

	if token != request.OneTimeToken {
		log.Printf("Invalid Token. \tReceived: %v \tWant: %v", request.OneTimeToken, token)
		http.Error(w, "Invalid token received", http.StatusUnauthorized)
		return
	}
	log.Printf("Successful OTC Match")

	// Get Licence details
	var licence utils.License
	fmt.Println("Retrieving Licence")

	if api.db == nil {
		log.Printf("Error: No pointer to db")
		http.Error(w, "Couldn't contact DB", http.StatusNotFound)
		return
	}

	query := "SELECT licence_type, machine_id, licence_key FROM licences WHERE user_id=$1"
	if err := api.db.QueryRow(query, request.UserID).Scan(&licence.LicenseType, &licence.MachineId, &licence.LicenseKey); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("error no rows matching. %v", err)
			http.Error(w, "user id not found in licence table", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("error in select statement: %v", err)
		return
	}

	// If machine id is new or nil update it
	if licence.MachineId == nil || *licence.MachineId != request.MachineID {
		// *licence.MachineId = request.MachineID
		_, err = api.db.Exec(`
			UPDATE licences
			SET machine_id = $1
			WHERE licence_key = $2`,
			request.MachineID, licence.LicenseKey)
		if err != nil {
			log.Printf("error updating machine_id %v", err.Error())
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	// Create JWT
	params := jwt.CreateJWTParams{
		UserId:     *request.UserID,
		MachineId:  licence.MachineId,
		LicenceKey: licence.LicenseKey,
		Plan:       licence.LicenseType,
	}

	jwt, err := jwt.CreateJWT(params)
	if err != nil {
		log.Printf("error encoding jwt %v", err.Error())
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data := map[string]string{"jwt": jwt}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("error encoding json %v", err.Error())
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	fmt.Println("JWT issued to :", *request.UserID)
}
