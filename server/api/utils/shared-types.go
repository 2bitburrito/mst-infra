package utils

type License struct {
	UserId      string   `json:"user_id"`
	MachineId   *string  `json:"machine_id"`
	LicenseType PlanType `json:"license_type"`
	LicenseKey  string   `json:"license_key"`
}
type JWT struct {
	Iss        string `json:"iss"`
	Sub        string `json:"sub"`
	MachineID  string `json:"machine"`
	Plan       string `json:"plan"`
	LicenceKey string `json:"licenceKey"`
	Exp        string `json:"exp"`
}

type PlanType string

const (
	PlanPaid  PlanType = "paid"
	PlanTrial PlanType = "trial"
	PlanBeta  PlanType = "beta"
)
