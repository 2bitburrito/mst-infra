package main

import (
	"fmt"
	"net/http"
)

type User struct {
	Email            string `json:"email"`
	HasLicense       bool   `json:"has_license"`
	NumberOfLicenses int    `json:"number_of_licenses"`
	FullName         string `json:"full_name"`
	Id               string `json:"Id"`
}

func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching User")
	id := r.PathValue("id")
	if len(id) < 1 {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}
  query := "SELECT email, has_license, number_of_licenses,  FROM users"
  db.Query("", args ...any)
	fmt.Println("userId:", id)
}
