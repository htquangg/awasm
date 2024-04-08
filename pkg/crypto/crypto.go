package crypto

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/schemas"

	"github.com/GoKillers/libsodium-go/cryptobox"
	generichash "github.com/GoKillers/libsodium-go/cryptogenerichash"
	"github.com/GoKillers/libsodium-go/cryptokdf"
	cryptosecretbox "github.com/GoKillers/libsodium-go/cryptosecretbox"
	"github.com/segmentfault/pacman/errors"
)

const (
	LOGIN_SUB_KEY_LENGTH      = 32
	LOGIN_SUB_KEY_ID          = 1
	LOGIN_SUB_KEY_CONTEXT     = "loginctx"
	LOGIN_SUB_KEY_BYTE_LENGTH = 16
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

	return schemas.EncryptionResult{Cipher: encryptedEmailBytes, Nonce: nonce, Key: encryptionKey}, nil
}

func GenerateKeyAndEncrypt(data string) (schemas.EncryptionResult, error) {
	encryptionKey, err := GenerateRandomBytes(cryptosecretbox.CryptoSecretBoxKeyBytes())
	if err != nil {
		return schemas.EncryptionResult{}, err
	}

	return Encrypt(data, encryptionKey)
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

func GetDecryptedToken(encryptedToken string, publicKey string, privateKey string) (string, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	encryptedTokenBytes, err := base64.StdEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	decryptedTokenBytes, errCode := cryptobox.CryptoBoxSealOpen(encryptedTokenBytes, publicKeyBytes, privateKeyBytes)
	if errCode != 0 {
		return "", errors.InternalServer(reason.UnknownError).WithMsg("Encryption token failed.")
	}

	return base64.URLEncoding.EncodeToString([]byte(decryptedTokenBytes)), nil
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

func GenerateKeyPair() (string, string, error) {
	skBytes, pkBytes, _ := cryptobox.CryptoBoxKeyPair()
	return base64.StdEncoding.EncodeToString(skBytes), base64.StdEncoding.EncodeToString(pkBytes), nil
}

func GenerateLoginSubKey(kek string) (string, error) {
	kekBytes, err := base64.StdEncoding.DecodeString(kek)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err)
	}

	kekSubKeyBytes, _ := cryptokdf.CryptoKdfDeriveFromKey(
		LOGIN_SUB_KEY_LENGTH,
		LOGIN_SUB_KEY_ID,
		LOGIN_SUB_KEY_CONTEXT,
		kekBytes,
	)
	// use first 16 bytes of generated kekSubKey as loginSubKey
	loginSubKeyBytes := kekSubKeyBytes[:16]

	return base64.StdEncoding.EncodeToString(loginSubKeyBytes), nil
}
