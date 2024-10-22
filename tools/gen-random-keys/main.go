package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/blake2b"

	"github.com/htquangg/awasm/pkg/converter"
	"github.com/htquangg/awasm/pkg/crypto"
)

func main() {
	keyBytes, err := crypto.GenerateRandomBytes(blake2b.Size256)
	if err != nil {
		log.Fatal(err)
	}
	key := converter.ToB64(keyBytes)

	hashBytes, err := crypto.GenerateRandomBytes(blake2b.Size)
	if err != nil {
		log.Fatal(err)
	}
	hash := converter.ToB64(hashBytes)

	apiKeySignatureHMACBytes, err := crypto.GenerateRandomBytes(blake2b.Size)
	if err != nil {
		log.Fatal(err)
	}
	apiKeySignatureHMAC := converter.ToURLB64(apiKeySignatureHMACBytes)

	apiKeyDatabaseHMACBytes, err := crypto.GenerateRandomBytes(blake2b.Size)
	if err != nil {
		log.Fatal(err)
	}
	apiKeyDatabaseHMAC := converter.ToURLB64(apiKeyDatabaseHMACBytes)

	cacheKeyHMACBytes, err := crypto.GenerateRandomBytes(blake2b.Size)
	if err != nil {
		log.Fatal(err)
	}
	cacheKeyHMAC := converter.ToB64(cacheKeyHMACBytes)

	jwtBytes, err := crypto.GenerateRandomBytes(blake2b.Size256)
	if err != nil {
		log.Fatal(err)
	}
	jwt := converter.ToURLB64(jwtBytes)

	fmt.Printf("key.encryption: %s\n", key)
	fmt.Printf("key.hash: %s\n", hash)
	fmt.Printf("key.api_key_signature_hmac: %s\n", apiKeySignatureHMAC)
	fmt.Printf("key.api_key_database_hmac: %s\n", apiKeyDatabaseHMAC)
	fmt.Printf("key.cache_key_hmac: %s\n", cacheKeyHMAC)
	fmt.Printf("jwt.secret: %s\n", jwt)
}
