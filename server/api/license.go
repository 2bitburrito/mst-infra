package main

import (
	"net/http"
)

type License struct {
	UserId      string   `json:"user_id"`
	MachineId   string   `json:"machine_id"`
	LicenseType PlanType `json:"license_type"`
	LicenseKey  string   `json:"license_key"`
}

type PlanType string

const (
	PlanPaid  PlanType = "paid"
	PlanTrial PlanType = "trial"
	PlanBeta  PlanType = "beta"
)

func postLicense(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not yet implemented", http.StatusNotFound)
}

func patchLicense(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not yet implemented", http.StatusNotFound)
}

func getLicense(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not yet implemented", http.StatusNotFound)
}

func checkLicense(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not yet implemented", http.StatusNotFound)
}
