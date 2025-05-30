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

func (api *API) checkLicense(w http.ResponseWriter, r *http.Request) {
	var jwtTokenString string
	var dbLicence queries.Licence

	err := json.NewDecoder(r.Body).Decode(&jwtTokenString)
	if err != nil {
		log.Printf("error decoding json in checkLoginCode %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
		// If db doesn't have machineID then instert:
		err := api.queries.ChangeMachineID(api.ctx, params)
		if err != nil {
			log.Printf("error adding the machine id to db %v", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	} else if claims.MachineID != dbLicence.MachineID.String {
		err := api.queries.RemoveMachineID(api.ctx, claims.LicenceKey)
		if err != nil {
			log.Printf("error adding the machine id to db %v", err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}
}
