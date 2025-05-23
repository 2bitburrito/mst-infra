package main

import (
	"crypto/ecdsa"

	"github.com/golang-jwt/jwt/v5"
)

func createJWT(plan PlanType, userId, machineId, licenceKey string) (string, error) {
	var (
		key *ecdsa.PrivateKey
		t   *jwt.Token
		s   string
	)

	key, err := loadPrivateKey("../private.pem")
	if err != nil {
		return "", err
	}

	t = jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"iss":        "Meta-Sound-Tools",
			"sub":        userId,
			"machine":    machineId,
			"plan":       plan,
			"licenceKey": licenceKey,
			"exp":        nil,
		})

	s, err = t.SignedString(key)
	if err != nil {
		return "", err
	}

	return s, nil
}
