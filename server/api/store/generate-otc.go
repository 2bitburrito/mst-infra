package store

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateOTC() string {
	arr := make([]byte, 6)
	rand.Read(arr)
	return hex.EncodeToString(arr)
}
