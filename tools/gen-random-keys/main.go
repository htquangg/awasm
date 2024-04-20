package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/htquangg/a-wasm/pkg/crypto"

	"golang.org/x/crypto/blake2b"
)

func main() {
	keyBytes, err := crypto.GenerateRandomBytes(blake2b.Size256)
	if err != nil {
		log.Fatal(err)
	}
	key := base64.StdEncoding.EncodeToString(keyBytes)

	hashBytes, err := crypto.GenerateRandomBytes(blake2b.Size)
	if err != nil {
		log.Fatal(err)
	}
	hash := base64.StdEncoding.EncodeToString(hashBytes)

	jwtBytes, err := crypto.GenerateRandomBytes(blake2b.Size256)
	if err != nil {
		log.Fatal(err)
	}
	jwt := base64.URLEncoding.EncodeToString(jwtBytes)

	fmt.Printf("key.encryption: %s\n", key)
	fmt.Printf("key.hash: %s\n", hash)
	fmt.Printf("jwt.secret: %s\n", jwt)
}
