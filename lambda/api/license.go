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
	PlanTrial PlanType = "trial"
	PlanPaid  PlanType = "paid"
	PlanBeta  PlanType = "beta"
)

func postLicense(w http.ResponseWriter, r *http.Request) {
}

func patchLicense(w http.ResponseWriter, r *http.Request) {
}

func getLicense(w http.ResponseWriter, r *http.Request) {
}

func checkLicense(w http.ResponseWriter, r *http.Request) {
}
