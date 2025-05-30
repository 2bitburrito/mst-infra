package utils

import (
	"fmt"
	"testing"
)

func TestLoadPublicKey(t *testing.T) {
	key, err := LoadPublicKey()
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
	}
	fmt.Println("Fetched Key:", key)
}
