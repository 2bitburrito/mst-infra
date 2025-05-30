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
	http.Error(w, "Method not yet implemented", http.StatusNotFound)
}

func (api *API) patchLicense(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not yet implemented", http.StatusNotFound)
}

func (api *API) getLicense(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not yet implemented", http.StatusNotFound)
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
		log.Printf("jwt invalid %v", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Get licence from db
	dbLicence, err = api.queries.GetLicence(api.ctx, claims.LicenceKey)
	if err != nil {
		log.Printf("error getting licence from db in checkLoginCode %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		err := api.queries.ChangeMachineID(api.ctx, params)
		if err != nil {
			log.Printf("error adding the machine id to db %v", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	} else if claims.MachineID != dbLicence.MachineID.String {
		log.Println("Machine Id doesn't match db - force logout!")
		err := api.queries.RemoveMachineID(api.ctx, claims.LicenceKey)
		if err != nil {
			log.Printf("error adding the machine id to db %v", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		responseBody := checkLicenceResponse{
			Message: "Machine ID doesn't match stored ID in database - Force logout",
			Action:  "logout",
		}
		dat, err := json.Marshal(responseBody)
		if err != nil {
			log.Printf("error marshalling json %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(dat)
	}

	w.WriteHeader(http.StatusOK)
	responseBody := checkLicenceResponse{
		Message: "License check successful",
		Action:  "logout",
	}
	dat, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("Error marshalling JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(dat)
}
