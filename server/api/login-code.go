package main

import (
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

	log.Printf("Retrieved OTC from store: %v", token)

	if token != request.OneTimeToken {
		returnJsonError(w, "Invalid token received", http.StatusUnauthorized)
		return
	}
	log.Printf("Successful OTC Match")

	log.Println("Retrieving Licence")

	// Get all licences of user:
	licences, err := api.queries.GetAllLicencesFromUserID(r.Context(), *request.UserID)
	if err != nil {
		returnJsonError(w, "error in select statement: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var newLicence database.Licence
	var oldestLicence database.Licence

	for _, licence := range licences {
		// First check whether licence is a plan or within expiry
		if !utils.LicenceIsValid(licence) {
			continue
		}
		// Track the oldest licence:
		if !oldestLicence.LastUsedAt.Valid {
			oldestLicence = licence
		}
		if oldestLicence.LastUsedAt.Time.After(licence.LastUsedAt.Time) {
			oldestLicence = licence
		}
		// If licence doesn't have a machine ID attached then attach it
		if !licence.MachineID.Valid {
			newLicence = licence
			newLicence.MachineID.String = request.MachineID
			err := api.queries.ChangeMachineIDAndJTI(r.Context(), database.ChangeMachineIDAndJTIParams{
				LicenceKey: newLicence.LicenceKey,
				MachineID:  newLicence.MachineID,
				Jti: uuid.NullUUID{
					Valid: false,
				},
			})
			if err != nil {
				returnJsonError(w, "error inserting new machine ID"+err.Error(), http.StatusInternalServerError)
				return
			}
			break
		} else if licence.MachineID.String == request.MachineID {
			newLicence = licence
			break
		}
	}

	if newLicence.UserID == uuid.Nil {
		// Defaulting back to oldest available licence
		newLicence = oldestLicence
		log.Printf("Reassigning license %s from machine %s to %s for user %s",
			oldestLicence.LicenceKey, oldestLicence.MachineID.String, request.MachineID, request.UserID.String())

		err := api.queries.ChangeMachineIDAndJTI(r.Context(), database.ChangeMachineIDAndJTIParams{
			LicenceKey: newLicence.LicenceKey,
			MachineID:  newLicence.MachineID,
			Jti: uuid.NullUUID{
				Valid: false,
			},
		})
		if err != nil {
			returnJsonError(w, "error updating machine ID for reassigned license"+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if len(newLicence.UserID) == 0 {
		returnJsonError(w, "No valid Licences Found", http.StatusUnauthorized)
		return
	}

	expiry := newLicence.Expiry.Time.Unix()
	jti := uuid.NullUUID{
		Valid: true,
		UUID:  uuid.New(),
	}

	// Create JWT
	params := jwt.Claims{
		UserID:     *request.UserID,
		MachineID:  newLicence.MachineID.String,
		LicenceKey: newLicence.LicenceKey,
		Plan:       utils.PlanType(newLicence.LicenceType.LicenceTypeEnum),
		Expiry:     expiry,
		JTI:        jti.UUID,
	}

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
