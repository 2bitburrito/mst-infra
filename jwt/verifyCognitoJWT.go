package jwt

import (
	"log"

	cognitoJwtVerify "github.com/jhosan7/cognito-jwt-verify"
)

func VerifyCognitoJWT(cognitoPoolId, userId, jwt string) (bool, error) {
	cognitoCfg := cognitoJwtVerify.Config{
		UserPoolId: cognitoPoolId,
		ClientId:   userId,
	}

	verify, err := cognitoJwtVerify.Create(cognitoCfg)
	if err != nil {
		return false, err
	}

	payload, err := verify.Verify(jwt)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return false, err
	}
	_, err = payload.GetSubject()
	if err != nil {
		log.Printf("Error: %s\n", err)
		return false, err
	}

	return true, nil
}

// NOTE:
// map[aud:3c2nrjdvdu92tmld1dhstg2bhu
// auth_time:1.747172729e+09
// cognito:username:e93919be-10a1-70c8-385b-c006812f9142
// email:hughpalmerproduction@gmail.com
// email_verified:true
// event_id:7d5c6d58-8b85-47a9-9e36-3323f6c5fd06
// exp:1.747282593e+09
// iat:1.747278993e+09
// iss:https://cognito-idp.us-west-1.amazonaws.com/us-west-1_bSqlVIEuH
// jti:2918c28d-7af5-4e39-9b93-9967810c0893
// name:Hugh Palmer
// origin_jti:80fbf13a-9674-4bf9-a06c-25c7238317ac
// sub:e93919be-10a1-70c8-385b-c006812f9142
// token_use:id]
