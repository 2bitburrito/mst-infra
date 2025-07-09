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
	"github.com/2bitburrito/mst-infra/server/api/licence-check"
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

func (api *API) checkLoginCodeAndCreateJWT(w http.ResponseWriter, r *http.Request) {
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
			returnJsonError(w, "error getting otc from store: "+err.Error(), http.StatusInternalServerError)
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

	if token != request.OneTimeToken {
		returnJsonError(w, "Invalid token received", http.StatusUnauthorized)
		return
	}

	// Get all licences of user:
	licences, err := api.queries.GetAllLicencesFromUserID(r.Context(), *request.UserID)
	if err != nil {
		returnJsonError(w, "error in select statement: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Pass licence to licence check to get valid licence:
	newLicence, err := licence.Check(request.MachineID, licences)
	if err != nil {
		returnJsonError(w, "Couldn't check licence validity "+err.Error(), http.StatusInternalServerError)
		return
	}
	if newLicence.UserID == uuid.Nil {
		returnJsonError(w, "No valid Licences Found", http.StatusUnauthorized)
		return
	}

	expiry := newLicence.Expiry.Time.Unix()
	jti := uuid.NullUUID{
		Valid: true,
		UUID:  uuid.New(),
	}
	err = api.queries.ChangeMachineIDAndJTI(r.Context(), database.ChangeMachineIDAndJTIParams{
		LicenceKey: newLicence.LicenceKey,
		MachineID:  sql.NullString{Valid: true, String: request.MachineID},
		Jti:        jti,
	})
	if err != nil {
		returnJsonError(w, "error inserting new machine ID"+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create JWT
	params := jwt.Claims{
		UserID:     *request.UserID,
		MachineID:  request.MachineID,
		LicenceKey: newLicence.LicenceKey,
		Plan:       utils.PlanType(newLicence.LicenceType.LicenceTypeEnum),
		Expiry:     expiry,
		JTI:        jti.UUID,
	}
	log.Println("NEW JWT PARAMS:")
	utils.PrintPretty(params)

	jwtToken, err := jwt.CreateJWT(params)
	if err != nil {
		returnJsonError(w, "internal error"+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Created New JWT for user: ", params.UserID)
	jwtJSON, _ := json.MarshalIndent(params, "", "  ")
	log.Println("Claims:", string(jwtJSON))

	// Insert JTI into licence row
	err = api.queries.ChangeJTI(r.Context(), database.ChangeJTIParams{
		Jti:        jti,
		LicenceKey: newLicence.LicenceKey,
	})
	if err != nil {
		returnJsonError(w, "error inserting JTI into row"+err.Error(), http.StatusInternalServerError)
	}

	log.Printf("Added jti: %v to db for licence key: %v", jti.UUID, newLicence.LicenceKey)

	data := map[string]string{"jwt": jwtToken}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		returnJsonError(w, "internal error encoding json "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("JWT issued to userID:", *request.UserID)
}
