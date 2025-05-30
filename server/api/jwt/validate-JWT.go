package jwt

import (
	"errors"
	"fmt"

	"github.com/2bitburrito/mst-infra/server/api/utils"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT(tokenString string) (*Claims, error) {
	key, err := utils.LoadPublicKey()
	if err != nil {
		return nil, fmt.Errorf("couldn't find public.pem file %v", err.Error())
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claim")
	}

	claimsStruct, err := mapClaimsToStruct(claims)
	if err != nil {
		return nil, fmt.Errorf("couldn't map claims to struct: %v", err)
	}
	return claimsStruct, nil
}
