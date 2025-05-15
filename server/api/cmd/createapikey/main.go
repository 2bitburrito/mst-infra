package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	const keyName = "API_KEY"

	bytesArr := make([]byte, 32)
	rand.Read(bytesArr)

	newKey := hex.EncodeToString(bytesArr)

	fmt.Printf("Generated key: \n%v\n", newKey)

	envMap, err := godotenv.Read()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	envMap[keyName] = newKey

	err = godotenv.Write(envMap, ".env")
	if err != nil {
		log.Fatal("error writing .env file")
	}
	fmt.Printf("Updated %s in .env\n", keyName)
}
