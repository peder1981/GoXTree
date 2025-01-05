package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
)

func generateUUIDV4() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return hex.EncodeToString(uuid), nil
}

func main() {
	uuid, err := generateUUIDV4()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated UUID: %s\n", uuid)
}
