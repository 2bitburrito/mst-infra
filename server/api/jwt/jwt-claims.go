package jwt

import (
	"encoding/json"

	"github.com/2bitburrito/mst-infra/server/api/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Issuer     string         `json:"iss,omitempty"`
	IssuedAt   int64          `json:"iat"`
	Expiry     int64          `json:"exp"`
	JTI        uuid.UUID      `json:"jti,omitempty"`
	UserID     uuid.UUID      `json:"sub"`
	MachineID  string         `json:"machine"`
	Plan       utils.PlanType `json:"plan"`
	LicenceKey string         `json:"licenceKey"`
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
