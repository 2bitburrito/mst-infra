package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	queries "github.com/2bitburrito/mst-infra/db/sqlc"
	"github.com/2bitburrito/mst-infra/jwt"
	"github.com/google/uuid"
)

type checkLicenceResponse struct {
	Message string `json:"message"`
	Action  string `json:"action"`
}

func (api *API) checkLicense(w http.ResponseWriter, r *http.Request) {
	var dbLicence queries.Licence

	jwtTokenString := r.Header.Get("Authorization")

	// Validate JWT
	claims, err := jwt.ValidateJWT(jwtTokenString)
	if err != nil {
		returnJsonError(w, "jwt invalid: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Get licence from db
	dbLicence, err = api.queries.GetLicence(r.Context(), claims.LicenceKey)
	if err != nil {
		returnJsonError(w, "error getting licence from db in checkLoginCode: "+err.Error(), http.StatusBadRequest)
		return
	}

	// If db doesn't have jit then insert it and the JIT:
	if !dbLicence.Jti.Valid {
		log.Println("Machine Id isn't set - setting in DB")
		err := api.queries.ChangeMachineIDAndJTI(r.Context(), queries.ChangeMachineIDAndJTIParams{
			LicenceKey: dbLicence.LicenceKey,
			MachineID: sql.NullString{
				String: claims.MachineID,
				Valid:  true,
			},
			Jti: uuid.NullUUID{
				UUID:  claims.JTI,
				Valid: true,
			},
		})
		if err != nil {
			returnJsonError(w, "error adding the machine id to db: "+err.Error(), http.StatusBadRequest)
			return
		}
		// Check the jti from the request matches the db
	} else if claims.JTI != dbLicence.Jti.UUID {
		// if they don't match we log user out
		w.WriteHeader(http.StatusUnauthorized)
		responseBody := checkLicenceResponse{
			Message: "Machine ID doesn't match stored ID in database - Force logout",
			Action:  "logout",
		}
		dat, err := json.Marshal(responseBody)
		if err != nil {
			returnJsonError(w, "error marshalling json"+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(dat)
		// We don't remove the old jwt as it's created and inserted on log in
		return
	}

	responseBody := checkLicenceResponse{
		Message: "License check successful",
		Action:  "null",
	}
	respondWithJSON(w, http.StatusOK, responseBody)

	if err := api.queries.UpdateLicenceUsedTime(context.Background(), dbLicence.LicenceKey); err != nil {
		log.Printf("Error while setting last_seen_at in licences for %s\n %v\n", dbLicence.LicenceKey, err)
	}
}

func (api *API) removeMachineFromLicence(w http.ResponseWriter, r *http.Request) {
	jwtTokenString := r.Header.Get("Authorization")

	log.Println("Received JWT:", jwtTokenString)

	// Validate JWT
	claims, err := jwt.ValidateJWT(jwtTokenString)
	if err != nil {
		returnJsonError(w, "jwt invalid: "+err.Error(), http.StatusUnauthorized, "There was an error while trying to validate jwt")
		return
	}
	err = api.queries.RemoveMachineID(r.Context(), claims.LicenceKey)
	if err != nil {
		returnJsonError(w, err.Error(), http.StatusInternalServerError, "Error while removing machine from licence row")
		return
	}
	w.WriteHeader(http.StatusOK)
}
