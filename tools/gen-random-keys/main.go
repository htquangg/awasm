package main

import (
	"encoding/base64"
	"fmt"
	"log"

	generichash "github.com/GoKillers/libsodium-go/cryptogenerichash"
	secretbox "github.com/GoKillers/libsodium-go/cryptosecretbox"
	"github.com/GoKillers/libsodium-go/sodium"
	"github.com/htquangg/a-wasm/pkg/crypto"
)

func main() {
	sodium.Init()

	keyBytes, err := crypto.GenerateRandomBytes(secretbox.CryptoSecretBoxKeyBytes())
	if err != nil {
		log.Fatal(err)
	}
	key := base64.StdEncoding.EncodeToString(keyBytes)

	hashBytes, err := crypto.GenerateRandomBytes(generichash.CryptoGenericHashBytesMax())
	if err != nil {
		log.Fatal(err)
	}
	hash := base64.StdEncoding.EncodeToString(hashBytes)

	jwtBytes, err := crypto.GenerateRandomBytes(secretbox.CryptoSecretBoxKeyBytes())
	if err != nil {
		log.Fatal(err)
	}
	jwt := base64.URLEncoding.EncodeToString(jwtBytes)

	fmt.Printf("key.encryption: %s\n", key)
	fmt.Printf("key.hash: %s\n", hash)
	fmt.Printf("jwt.secret: %s\n", jwt)
}
