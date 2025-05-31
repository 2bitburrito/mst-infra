package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	database "github.com/2bitburrito/mst-infra/db/sqlc"
	"github.com/google/uuid"
)

type User struct {
	Id               string `json:"id"`
	Email            string `json:"email"`
	HasLicense       bool   `json:"has_license"`
	NumberOfLicenses int    `json:"number_of_licenses"`
	FullName         string `json:"full_name"`
	JWT              string `json:"jwt"`
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(data, &request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := request.Id

	log.Println("Received user ID:", id)

	if len(id) < 1 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	query := "SELECT email, has_license, number_of_licenses, id FROM users WHERE id=$1"
	if err := api.db.QueryRow(query, id).Scan(&user.Email, &user.HasLicense, &user.NumberOfLicenses, &user.Id); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("error no rows matching. %v", err)
			http.Error(w, "user id not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("error in select statement: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("error in select statement: %v", err)
		http.Error(w, "failed to encode user data to json", http.StatusInternalServerError)
		return
	}
}

func (api *API) patchUser(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not yet implemented", http.StatusNotFound)
}

func (api *API) postUser(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not yet implemented", http.StatusNotFound)
}

func (api *API) postCognitoUser(w http.ResponseWriter, r *http.Request) {
	var cognitoUser CognitoUser

	log.Println("Recieved Cognito Request for:", cognitoUser.Email)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(data, &cognitoUser); err != nil {
		log.Println("error unmarshalling json", err)
	}

	nonNullStr := sql.NullString{
		String: cognitoUser.Email,
		Valid:  true,
	}
	email, err := api.queries.GetBetaEmail(api.ctx, nonNullStr)
	if err != nil {
		log.Printf("error in while getting beta email: %v", err)
		http.Error(w, "error retrieving email from beta list", http.StatusInternalServerError)
	}
	if email.Valid {
		// this means we have a beta licence
		// So we update the userid correctly
		api.queries.UpdateUserId(api.ctx, database.UpdateUserIdParams{
			ID:    cognitoUser.Sub,
			Email: cognitoUser.Email,
		})
		w.WriteHeader(http.StatusOK)
		log.Println("Updated user ID for Beta User: ", cognitoUser.Email)
		return
	}

	args := database.InsertUserParams{
		ID:                 cognitoUser.Sub,
		Email:              cognitoUser.Email,
		FullName:           cognitoUser.Name,
		HasLicense:         false,
		SubscribedToEmails: false,
	}
	if err := api.queries.InsertUser(api.ctx, args); err != nil {
		log.Printf("error in while writing cognito user to db: %v", err)
		http.Error(w, "error in while writing cognito user to db", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (api *API) deleteUser(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not yet implemented", http.StatusNotFound)
}
