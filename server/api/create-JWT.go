package main

import (
	"crypto/ecdsa"

	"github.com/golang-jwt/jwt/v5"
)

func createJWT(plan PlanType, userId, machineId, licenceKey string) (string, error) {
	var (
		key       *ecdsa.PrivateKey
		token     *jwt.Token
		jwtString string
	)

	key, err := loadPrivateKey("../private.pem")
	if err != nil {
		return "", err
	}
	// TODO: Fetch the correct licence and insure the correct exp is set

	token = jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"iss":        "Meta-Sound-Tools",
			"sub":        userId,
			"machine":    machineId,
			"plan":       plan,
			"licenceKey": licenceKey,
			"exp":        nil,
		})

	jwtString, err = token.SignedString(key)
	if err != nil {
		return "", err
	}

	return jwtString, nil
}
