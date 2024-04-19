package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"io"

	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/pkg/converter"

	"github.com/segmentfault/pacman/errors"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	LOGIN_SUB_KEY_INFO = "loginInfo"
)

func Encrypt(data string, encryptionKey []byte) (schemas.EncryptionResult, error) {
	nonce, err := GenerateRandomBytes(24)
	if err != nil {
		return schemas.EncryptionResult{}, errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}

	encryptedEmailBytes := secretbox.Seal(nil, []byte(data), (*[24]byte)(nonce), (*[32]byte)(encryptionKey))

	return schemas.EncryptionResult{Cipher: encryptedEmailBytes, Nonce: nonce, Key: encryptionKey}, nil
}

func GenerateKeyAndEncrypt(data string) (schemas.EncryptionResult, error) {
	encryptionKey, err := GenerateRandomBytes(32)
	if err != nil {
		return schemas.EncryptionResult{}, err
	}

	return Encrypt(data, encryptionKey)
}

func Decrypt(cipher []byte, key []byte, nonce []byte) (string, error) {
	decryptedBytes, ok := secretbox.Open(nil, cipher[:], (*[24]byte)(nonce), (*[32]byte)(key))
	if !ok {
		return "", errors.InternalServer(reason.UnknownError).WithMsg("Decryption failed.").WithStack()
	}

	return string(decryptedBytes), nil
}

func GetHash(data string, hashKey []byte) (string, error) {
	h, err := blake2b.New256(hashKey)
	if err != nil {
		return "", err
	}
	_, err = h.Write([]byte(data))
	if err != nil {
		return "", err
	}
	return converter.ToB64(h.Sum(nil)), nil
}

func GetEncryptedToken(token string, publicKey string) (string, error) {
	publicKeyBytes, err := converter.FromB64(publicKey)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	tokenBytes, err := converter.FromURLB64(token)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	encryptedTokenBytes, err := box.SealAnonymous(nil, tokenBytes, (*[32]byte)(publicKeyBytes), nil)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithMsg("Encryption token failed.")
	}

	return converter.ToB64(encryptedTokenBytes), nil
}

func GetDecryptedToken(encryptedToken string, publicKey string, privateKey string) (string, error) {
	publicKeyBytes, err := converter.FromB64(publicKey)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	privateKeyBytes, err := converter.FromB64(privateKey)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	encryptedTokenBytes, err := converter.FromB64(encryptedToken)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	decryptedTokenBytes, ok := box.OpenAnonymous(
		nil,
		encryptedTokenBytes,
		(*[32]byte)(publicKeyBytes),
		(*[32]byte)(privateKeyBytes),
	)
	if !ok {
		return "", errors.InternalServer(reason.UnknownError).WithMsg("Encryption token failed.")
	}

	return converter.ToURLB64([]byte(decryptedTokenBytes)), nil
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
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}

	return converter.ToURLB64(b), nil
}

func GenerateKeyPair() (string, string, error) {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	return converter.ToB64((*priv)[:]), converter.ToB64((*pub)[:]), nil
}

func GenerateLoginSubKey(kek string) (string, error) {
	kekBytes, err := converter.FromB64(kek)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}

	loginSubKeyBytes := make([]byte, 16)
	kdf := hkdf.New(sha256.New, kekBytes, nil, []byte(LOGIN_SUB_KEY_INFO))
	if _, err := io.ReadFull(kdf, loginSubKeyBytes); err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}

	return converter.ToB64(loginSubKeyBytes), nil
}
