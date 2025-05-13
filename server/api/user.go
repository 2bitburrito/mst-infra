package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Id               string `json:"id"`
	Email            string `json:"email"`
	HasLicense       bool   `json:"has_license"`
	NumberOfLicenses int    `json:"number_of_licenses"`
	FullName         string `json:"full_name"`
}

func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching User")
	id := r.PathValue("id")
	fmt.Println("userId:", id)
	if len(id) < 1 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var user User
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

	fmt.Println("userId:", id)
}

func patchUser(w http.ResponseWriter, r *http.Request)  {}
func postUser(w http.ResponseWriter, r *http.Request)   {}
func deleteUser(w http.ResponseWriter, r *http.Request) {}
