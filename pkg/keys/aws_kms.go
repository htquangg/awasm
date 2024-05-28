package keys

import (
	"context"
	"crypto"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/lstoll/awskms"
)

var _ KeyManager = (*AWSKMS)(nil)

type AWSKMS struct {
	svc *kms.Client
}

func NewAWSKMS(ctx context.Context) (KeyManager, error) {
	defaultConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS KMS config: %w", err)
	}

	svc := kms.NewFromConfig(defaultConfig)

	return &AWSKMS{
		svc: svc,
	}, nil
}

func (s *AWSKMS) NewSigner(ctx context.Context, keyID string) (crypto.Signer, error) {
	return awskms.NewSigner(ctx, s.svc, keyID)
}

func (s *AWSKMS) Encrypt(
	ctx context.Context,
	keyID string,
	plaintext []byte,
	aad []byte,
) ([]byte, error) {
	input := kms.EncryptInput{
		KeyId: &keyID,
		EncryptionContext: map[string]string{
			"aad": base64.StdEncoding.EncodeToString(aad),
		},
		Plaintext: plaintext,
	}
	output, err := s.svc.Encrypt(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt: %w", err)
	}

	return output.CiphertextBlob, nil
}

func (s *AWSKMS) Decrypt(
	ctx context.Context,
	keyID string,
	ciphertext []byte,
	aad []byte,
) ([]byte, error) {
	input := kms.DecryptInput{
		KeyId: &keyID,
		EncryptionContext: map[string]string{
			"aad": base64.StdEncoding.EncodeToString(aad),
		},
		CiphertextBlob: ciphertext,
	}
	output, err := s.svc.Decrypt(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return output.Plaintext, nil
}
