package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	database "github.com/2bitburrito/mst-infra/db/sqlc"
	"github.com/2bitburrito/mst-infra/server/api/config"
	"github.com/2bitburrito/mst-infra/server/api/jwt"
	"github.com/2bitburrito/mst-infra/server/api/utils"
	"github.com/google/uuid"
)

type CreateLoginCodeRequest struct {
	UserID string `json:"id"`
	JWT    string `json:"jwt"`
}

type LoginCodeRequest struct {
	UserID       *uuid.UUID `json:"userId,omitempty"`
	OneTimeToken string     `json:"otc"`
	MachineID    string     `json:"machineId"`
}

func (api *API) createLoginCode(w http.ResponseWriter, r *http.Request) {
	cfg, _ := config.LoadConfig()
	var user User

	log.Println("Creating Login Code")

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		returnJsonError(w, "Couldn't decode json"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Verify the jwt with cognito
	verified, err := jwt.VerifyCognitoJWT(cfg.CognitoPoolID, user.Id.String(), user.JWT)
	if err != nil || !verified {
		returnJsonError(w, "error in Jwt Verify "+err.Error(), http.StatusNotFound)
		return
	}

	otc := api.verificationStore.New(user.Id)
	returnObj := map[string]string{
		"userId": user.Id.String(),
		"otc":    otc,
	}

	returnData, err := json.Marshal(returnObj)
	if err != nil {
		returnJsonError(w, "error marshalling otc: "+err.Error(), http.StatusNotFound)
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
		returnJsonError(w, "error decoding json in checkLoginCode %v"+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var token string

	if request.UserID == nil {
		// Get Token from store from OTC
		var userId uuid.UUID
		userId, token, err = api.verificationStore.GetFromOTC(request.OneTimeToken)
		if err != nil {
			returnJsonError(w, "error getting otc from store"+err.Error(), http.StatusInternalServerError)
			return
		}
		request.UserID = &userId
	} else {
		// Get Token from store matching userID:
		token, err = api.verificationStore.Get(*request.UserID)
		if err != nil {
			returnJsonError(w, "error getting otc from store"+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	log.Printf("Retrieved OTC from store: %v", token)

	if token != request.OneTimeToken {
		returnJsonError(w, "Invalid token received", http.StatusUnauthorized)
		return
	}
	log.Printf("Successful OTC Match")

	// Get Licence details
	var licence utils.License
	log.Println("Retrieving Licence")

	if api.db == nil {
		returnJsonError(w, "Couldn't contact DB. DB is nil", http.StatusNotFound)
		return
	}

	var licenceRow database.AddTrialLicenceRow
	query := "SELECT licence_type, machine_id, licence_key, expiry FROM licences WHERE user_id=$1"
	if err := api.db.QueryRow(query, *request.UserID).Scan(&licence.LicenseType, &licence.MachineId, &licence.LicenseKey); err != nil {
		if err == sql.ErrNoRows {
			// Create a trial licence
			args := database.AddTrialLicenceParams{
				UserID:    *request.UserID,
				MachineID: sql.NullString{String: request.MachineID, Valid: true},
			}
			licenceRow, err = api.queries.AddTrialLicence(api.ctx, args)
			if err != nil {
				returnJsonError(w, "error while adding trial user to table: "+err.Error(), http.StatusBadRequest)
				return
			}
			licence.LicenseKey = licenceRow.LicenceKey
		} else {
			returnJsonError(w, "error in select statement:"+err.Error(), http.StatusInternalServerError)
			return
		}
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
			returnJsonError(w, "error updating machine_id: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Create JWT
	params := jwt.CreateJWTParams{
		UserId:     *request.UserID,
		MachineId:  licence.MachineId,
		LicenceKey: licence.LicenseKey,
		Plan:       licence.LicenseType,
		Expiry:     licenceRow.Expiry,
	}

	jwtToken, err := jwt.CreateJWT(params)
	if err != nil {
		returnJsonError(w, "internal error"+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{"jwt": jwtToken}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		returnJsonError(w, "internal error encoding json "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("JWT issued to :", *request.UserID)
}
