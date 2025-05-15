package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
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

// 	data, err = io.ReadAll(res.Body)
// 	if err != nil {
// 		return exploreRes, err
// 	}
// 	c.Cache.Add(path, data)
// }
// if err = json.Unmarshal(data, &exploreRes); err != nil {
// 	return exploreRes, err

func getUser(w http.ResponseWriter, r *http.Request) {
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
	if err := db.QueryRow(query, id).Scan(&user.Email, &user.HasLicense, &user.NumberOfLicenses, &user.Id); err != nil {
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

func patchUser(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not yet implemented", http.StatusNotFound)
}

func postUser(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not yet implemented", http.StatusNotFound)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not yet implemented", http.StatusNotFound)
}
