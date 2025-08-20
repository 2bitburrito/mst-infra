package utils

import (
	"fmt"
	"testing"
)

func TestLoadPublicKey(t *testing.T) {
	publicKeyPath := "../public.pem"
	key, err := LoadPublicKey(publicKeyPath)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
	}
	fmt.Println("Fetched Key:", key)
}
