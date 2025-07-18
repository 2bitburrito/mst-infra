package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/2bitburrito/mst-infra/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func mockLoadPrivateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func TestCreateJWTAndUnmarshalClaims(t *testing.T) {
	origLoadPrivateKey := utils.LoadPrivateKey
	// Replace with mock
	utils.LoadPrivateKey = mockLoadPrivateKey
	// Restore after test
	defer func() { utils.LoadPrivateKey = origLoadPrivateKey }()

	userID := uuid.New()
	machineID := "test-machine"

	now := time.Now().Unix()
	params := Claims{
		UserID:     userID,
		MachineID:  machineID,
		LicenceKey: "test-licence",
		Plan:       "beta",
		IssuedAt:   now,
		Expiry:     time.Now().Add(time.Hour).Unix(),
	}

	tokenString, err := CreateJWT(params)
	require.NoError(t, err)

	parsedToken, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	require.NoError(t, err)

	mapClaims := parsedToken.Claims.(jwt.MapClaims)
	claims, err := mapClaimsToStruct(mapClaims)
	claimsJSON, _ := json.MarshalIndent(claims, "", "  ")
	fmt.Println("Claims:", string(claimsJSON))
	require.NoError(t, err)

	require.Equal(t, userID.String(), claims.UserID.String())
	require.Equal(t, machineID, claims.MachineID)
	require.Equal(t, utils.PlanType("beta"), claims.Plan)
	require.Equal(t, "test-licence", claims.LicenceKey)
	require.Equal(t, params.Expiry, claims.Expiry)
	require.Equal(t, now, claims.IssuedAt)

	expectedFields := []string{
		"iss", "sub", "jti", "iat", "exp",
		"machine", "plan", "licenceKey",
	}
	for _, field := range expectedFields {
		require.Contains(t, mapClaims, field, "JWT missing field: %s", field)
	}
}
