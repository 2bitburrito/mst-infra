package jwt

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID     string `json:"sub"`
	MachineID  string `json:"machine"`
	Plan       string `json:"plan"`
	LicenceKey string `json:"licenceKey"`
	Expiry     string `json:"exp"`
}

func mapClaimsToStruct(claims jwt.MapClaims) (*Claims, error) {
	jsonBytes, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}

	var result Claims
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
