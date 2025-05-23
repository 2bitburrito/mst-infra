package main

import (
	"testing"
)

func TestLoadPrivateKey(t *testing.T) {
	_, err := loadPrivateKey("../private.pem")
	if err != nil {
		t.Error("error while loading priv key:", err.Error())
	}
}
