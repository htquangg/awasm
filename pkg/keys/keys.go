package keys

import (
	"context"
	"crypto"
)

type KeyManager interface {
	NewSigner(ctx context.Context, keyID string) (crypto.Signer, error)
	Encrypt(ctx context.Context, keyID string, plaintext []byte, aad []byte) ([]byte, error)
	Decrypt(ctx context.Context, keyID string, ciphertext []byte, aad []byte) ([]byte, error)
}
