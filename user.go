package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

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

type CognitoUser struct {
	Sub                uuid.UUID `json:"sub"`
	Email              string    `json:"email"`
	Name               string    `json:"name"`
	ConfirmationStatus string    `json:"user_status"`
}

type returnUserPayload struct {
	Email              string    `json:"email"`
	CreatedAt          time.Time `json:"createdAt"`
	NumberOfLicences   int32     `json:"numberOfLicences"`
	SubscribedToEmails bool      `json:"subscribedToEmails"`
	FullName           string    `json:"fullName"`
	ID                 uuid.UUID `json:"id"`
}

func (api *API) getUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if len(id) == 0 {
		returnJsonError(w, "Invalid id: "+id, http.StatusBadRequest)
		return
	}
	UUIDid := uuid.MustParse(id)

	user, err := api.queries.GetUser(r.Context(), UUIDid)
	if err != nil {
		returnJsonError(w, "Error while getting user in db "+err.Error(), 500)
		return
	}
	userPayload := returnUserPayload{
		ID:                 user.ID,
		Email:              user.Email,
		CreatedAt:          user.CreatedAt.Time,
		NumberOfLicences:   user.NumberOfLicences,
		SubscribedToEmails: user.SubscribedToEmails,
		FullName:           user.FullName,
	}

	respondWithJSON(w, http.StatusOK, userPayload)
}

// When user is created in cognito we either assign them to a trial licence or Beta licence
func (api *API) postCognitoUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
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
	log.Printf("Time to unmarshal: %v\n", time.Since(start))

	nonNullEmail := sql.NullString{
		String: cognitoUser.Email,
		Valid:  true,
	}

	betaLicence, err := api.queries.GetBetaEmail(r.Context(), nonNullEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			betaLicence.Email.Valid = false
		} else {
			returnJsonError(w, "error retrieving email from beta list: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	log.Printf("Time to check Beta Licences: %v\n", time.Since(start))
	args := database.InsertUserParams{
		ID:                 cognitoUser.Sub,
		Email:              cognitoUser.Email,
		FullName:           cognitoUser.Name,
		SubscribedToEmails: true,
	}
	log.Println("Inserting user: ", args)
	if err := api.queries.InsertUser(r.Context(), args); err != nil {
		returnJsonError(w, "error in while writing cognito user to db: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Time to insert user: %v\n", time.Since(start))

	if betaLicence.Email.Valid {
		// If user is in beta list:
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
	} else {
		// Add a trial licence
		licenceRowArgs := database.AddTrialLicenceParams{
			UserID: cognitoUser.Sub,
			MachineID: sql.NullString{
				Valid: false,
			},
		}
		licenceRow, err := api.queries.AddTrialLicence(r.Context(), licenceRowArgs)
		if err != nil {
			returnJsonError(w, "error while adding trial user to table: "+err.Error(), http.StatusBadRequest)
			return
		}
		log.Println("Added new Trial Licence: ", licenceRow)
		log.Printf("Time to insert Trial Licence: %v\n", time.Since(start))
	}
	w.WriteHeader(http.StatusOK)
}

func (api *API) deleteUser(w http.ResponseWriter, r *http.Request) {
	returnJsonError(w, "Method not yet implemented", http.StatusNotFound)
}

func (api *API) checkUserIsBeta(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")

	_, err := api.queries.GetBetaEmail(r.Context(), sql.NullString{
		Valid:  true,
		String: email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			returnJsonError(w, "Email not enrolled in Beta Program", http.StatusNotFound)
			return
		} else {
			returnJsonError(w, "Error while fetching user in beta: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
