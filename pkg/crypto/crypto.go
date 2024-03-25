package crypto

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/schemas"

	"github.com/GoKillers/libsodium-go/cryptobox"
	generichash "github.com/GoKillers/libsodium-go/cryptogenerichash"
	cryptosecretbox "github.com/GoKillers/libsodium-go/cryptosecretbox"
	"github.com/segmentfault/pacman/errors"
)

func Encrypt(data string, encryptionKey []byte) (schemas.EncryptionResult, error) {
	nonce, err := GenerateRandomBytes(cryptosecretbox.CryptoSecretBoxNonceBytes())
	if err != nil {
		return schemas.EncryptionResult{}, errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	encryptedEmailBytes, errCode := cryptosecretbox.CryptoSecretBoxEasy([]byte(data), nonce, encryptionKey)
	if errCode != 0 {
		return schemas.EncryptionResult{}, errors.InternalServer(reason.UnknownError).WithMsg("Encryption failed.")
	}

	return schemas.EncryptionResult{Cipher: encryptedEmailBytes, Nonce: nonce}, nil
}

func Decrypt(cipher []byte, key []byte, nonce []byte) (string, error) {
	decryptedBytes, err := cryptosecretbox.CryptoSecretBoxOpenEasy(cipher, nonce, key)
	if err != 0 {
		return "", errors.InternalServer(reason.UnknownError).WithMsg("Decryption failed.")
	}

	return string(decryptedBytes), nil
}

func GetHash(data string, hashKey []byte) (string, error) {
	dataHashBytes, err := generichash.CryptoGenericHash(generichash.CryptoGenericHashBytes(), []byte(data), hashKey)
	if err != 0 {
		return "", errors.InternalServer(reason.UnknownError).WithMsg("Hash failed.")
	}

	return base64.StdEncoding.EncodeToString(dataHashBytes), nil
}

func GetEncryptedToken(token string, publicKey string) (string, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	tokenBytes, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	encryptedTokenBytes, errCode := cryptobox.CryptoBoxSeal(tokenBytes, publicKeyBytes)
	if errCode != 0 {
		return "", errors.InternalServer(reason.UnknownError).WithMsg("Encryption token failed.")
	}

	return base64.StdEncoding.EncodeToString(encryptedTokenBytes), nil
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, errors.InternalServer(reason.UnknownError).
			WithError(err).
			WithStack().
			WithMsg("Generation random bytes failed.")
	}

	return b, nil
}

func GenerateURLSafeRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
