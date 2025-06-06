package jwt

import (
	"crypto/ecdsa"
	"log"
	"time"

	"github.com/2bitburrito/mst-infra/server/api/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateJWT(params Claims) (string, error) {
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
			"sub":        params.UserID,
			"machine":    params.MachineID,
			"plan":       params.Plan,
			"licenceKey": params.LicenceKey,
			"exp":        params.Expiry,
			"jti":        uuid.NewString(),
			"iat":        time.Now().Unix(),
		})

	jwtString, err = token.SignedString(key)
	if err != nil {
		return "", err
	}

	return jwtString, nil
}
