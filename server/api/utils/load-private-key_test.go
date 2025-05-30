package utils

import (
	"testing"
)

func TestLoadPrivateKey(t *testing.T) {
	_, err := LoadPrivateKey()
	if err != nil {
		t.Error("error while loading priv key:", err.Error())
	}
}
