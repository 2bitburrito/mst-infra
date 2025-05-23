package main

import (
	"fmt"
	"testing"
)

func TestCreateJWT(t *testing.T) {
	jwt, err := createJWT("perpetual", "1", "mac_02394", "2098349djklj")
	if err != nil {
		t.Errorf("couldn't create jwt %v", err.Error())
	}
	fmt.Println("JWT:", jwt)
}
