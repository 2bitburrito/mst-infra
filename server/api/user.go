package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	database "github.com/2bitburrito/mst-infra/db/sqlc"
	"github.com/google/uuid"
)

type User struct {
	Id               uuid.UUID `json:"id"`
	Email            string    `json:"email"`
	HasLicense       bool      `json:"has_license"`
	NumberOfLicenses int       `json:"number_of_licenses"`
	FullName         string    `json:"full_name"`
	JWT              string    `json:"jwt"`
}
type receivedUserRequest struct {
	Id string `json:"id"`
}

type CognitoUser struct {
	Sub                uuid.UUID `json:"sub"`
	Email              string    `json:"email"`
	Name               string    `json:"name"`
	ConfirmationStatus string    `json:"user_status"`
}

func (api *API) getUser(w http.ResponseWriter, r *http.Request) {
	var user User
	var request receivedUserRequest

	log.Println("Fetching User")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		returnJsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(data, &request); err != nil {
		returnJsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := request.Id

	log.Println("Received user ID:", id)

	if len(id) < 1 {
		returnJsonError(w, "Invalid id: "+id, http.StatusBadRequest)
		return
	}

	query := "SELECT email, has_license, number_of_licenses, id FROM users WHERE id=$1"
	if err := api.db.QueryRow(query, id).Scan(&user.Email, &user.HasLicense, &user.NumberOfLicenses, &user.Id); err != nil {
		if err == sql.ErrNoRows {
			returnJsonError(w, "user id not found: "+err.Error(), http.StatusNotFound)
			return
		}
		returnJsonError(w, "error in sql statement: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		returnJsonError(w, "failed to encode user data to json: "+err.Error(), http.StatusNotFound)
		return
	}
}

func (api *API) patchUser(w http.ResponseWriter, r *http.Request) {
	returnJsonError(w, "Method not yet implemented", http.StatusNotFound)
}

func (api *API) postUser(w http.ResponseWriter, r *http.Request) {
	returnJsonError(w, "Method not yet implemented", http.StatusNotFound)
}

func (api *API) postCognitoUser(w http.ResponseWriter, r *http.Request) {
	var cognitoUser CognitoUser

	data, err := io.ReadAll(r.Body)
	if err != nil {
		returnJsonError(w, "error reading body json: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(data, &cognitoUser); err != nil {
		returnJsonError(w, "error unmarshalling json: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Recieved Cognito Request for:", cognitoUser.Email)

	nonNullStr := sql.NullString{
		String: cognitoUser.Email,
		Valid:  true,
	}
	betaLicence, err := api.queries.GetBetaEmail(r.Context(), nonNullStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			betaLicence.Email.Valid = false
		} else {
			returnJsonError(w, "error retrieving email from beta list: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	args := database.InsertUserParams{
		ID:                 cognitoUser.Sub,
		Email:              cognitoUser.Email,
		FullName:           cognitoUser.Name,
		HasLicense:         false,
		SubscribedToEmails: false,
	}

	if betaLicence.Email.Valid {
		// If user is in beta list:
		args.HasLicense = true
		emailNonNull := sql.NullString{
			Valid:  true,
			String: args.Email,
		}
		err := api.queries.SetBetaRowToSeen(r.Context(), emailNonNull)
		if err != nil {
			log.Println("error setting beta row to seen: ", err.Error())
		}

		// add beta licence:
		_, err = api.queries.AddBetaLicence(r.Context(), args.ID)
		if err != nil {
			returnJsonError(w, "error while setting new beta licence:"+err.Error(), http.StatusInternalServerError)
		}
	}
	log.Println("Inserting user: ", args)
	if err := api.queries.InsertUser(r.Context(), args); err != nil {
		log.Println("ERROR NOW: ", err.Error())
		returnJsonError(w, "error in while writing cognito user to db: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *API) deleteUser(w http.ResponseWriter, r *http.Request) {
	returnJsonError(w, "Method not yet implemented", http.StatusNotFound)
}
