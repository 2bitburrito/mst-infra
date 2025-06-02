package jwt

import (
	"crypto/ecdsa"
	"log"

	"github.com/2bitburrito/mst-infra/server/api/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CreateJWTParams struct {
	UserId     uuid.UUID
	MachineId  *string
	LicenceKey string
	Plan       utils.PlanType
	Expiry     int64
}

func CreateJWT(params CreateJWTParams) (string, error) {
	var (
		key       *ecdsa.PrivateKey
		token     *jwt.Token
		jwtString string
	)

	key, err := utils.LoadPrivateKey()
	if err != nil {
		return "", err
	}

	log.Println("Plan Type:", params.Plan)

	token = jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"iss":        "Meta-Sound-Tools",
			"sub":        params.UserId,
			"machine":    params.MachineId,
			"plan":       params.Plan,
			"licenceKey": params.LicenceKey,
			"exp":        params.Expiry,
		})

	jwtString, err = token.SignedString(key)
	if err != nil {
		return "", err
	}

	return jwtString, nil
}
