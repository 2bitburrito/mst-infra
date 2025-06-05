package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	queries "github.com/2bitburrito/mst-infra/db/sqlc"
	"github.com/2bitburrito/mst-infra/server/api/jwt"
)

func (api *API) postLicense(w http.ResponseWriter, r *http.Request) {
	returnJsonError(w, "Method not yet implemented", http.StatusNotFound)
}

func (api *API) patchLicense(w http.ResponseWriter, r *http.Request) {
	returnJsonError(w, "Method not yet implemented", http.StatusNotFound)
}

func (api *API) getLicense(w http.ResponseWriter, r *http.Request) {
	returnJsonError(w, "Method not yet implemented", http.StatusNotFound)
}

type checkLicenceResponse struct {
	Message string `json:"message"`
	Action  string `json:"action"`
}

func (api *API) checkLicense(w http.ResponseWriter, r *http.Request) {
	var dbLicence queries.Licence

	jwtTokenString := r.Header.Get("Authorization")

	log.Println("Received JWT:", jwtTokenString)

	// Validate JWT
	claims, err := jwt.ValidateJWT(jwtTokenString)
	if err != nil {
		returnJsonError(w, "jwt invalid", http.StatusUnauthorized)
		return
	}

	// Get licence from db
	dbLicence, err = api.queries.GetLicence(r.Context(), claims.LicenceKey)
	if err != nil {
		returnJsonError(w, "error getting licence from db in checkLoginCode: "+err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Retrieved licence for userid:%s / licenceKey:%s", dbLicence.UserID, dbLicence.LicenceKey)

	params := queries.ChangeMachineIDParams{
		LicenceKey: dbLicence.LicenceKey,
		MachineID: sql.NullString{
			String: claims.MachineID,
			Valid:  true,
		},
	}

	if !dbLicence.MachineID.Valid {
		log.Println("Machine Id isn't set - setting in DB")
		// If db doesn't have machineID then instert:
		err := api.queries.ChangeMachineID(r.Context(), params)
		if err != nil {
			returnJsonError(w, "error adding the machine id to db: "+err.Error(), http.StatusBadRequest)
			return
		}

	} else if claims.MachineID != dbLicence.MachineID.String {
		// Licence is using a new machine
		log.Println("Machine Id doesn't match db - force logout!")
		err := api.queries.RemoveMachineID(r.Context(), claims.LicenceKey)
		if err != nil {
			returnJsonError(w, "error adding the machine id to db: "+err.Error(), http.StatusBadRequest)
			return
		}
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
	}

	w.WriteHeader(http.StatusOK)
	responseBody := checkLicenceResponse{
		Message: "License check successful",
		Action:  "null",
	}
	dat, err := json.Marshal(responseBody)
	if err != nil {
		returnJsonError(w, "Error marshalling JSON response: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(dat)
}
