package api_key

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/pkg/converter"
	"github.com/htquangg/a-wasm/pkg/crypto"
	"github.com/htquangg/a-wasm/pkg/uid"
)

const (
	apiKeyBytes = 64 // 64 bytes is 86 chararacters in non-padded base64.
)

type ApiKeyRepo interface {
	AddApiKey(ctx context.Context, apiKey *entities.ApiKey) error
	GetApiKeyByKey(ctx context.Context, key string) (*entities.ApiKey, bool, error)
}

type ApiKeyService struct {
	cfg        *config.Config
	apiKeyRepo ApiKeyRepo
}

func NewApiKeyService(cfg *config.Config, apiKeyRepo ApiKeyRepo) *ApiKeyService {
	return &ApiKeyService{
		cfg:        cfg,
		apiKeyRepo: apiKeyRepo,
	}
}

func (s *ApiKeyService) AddApiKey(
	ctx context.Context,
	req *schemas.AddApiKeyReq,
) (*schemas.AddApiKeyResp, error) {
	fullApiKey, err := s.generateApiKey(req.UserID)
	if err != nil {
		return nil, errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}

	parts := strings.SplitN(fullApiKey, ".", 3)
	if len(parts) != 3 {
		return nil, errors.InternalServer(reason.ApiKeyInvalid).WithStack()
	}
	key := parts[0]

	hmacedKey, err := s.generateApiKeyHMAC(key)
	if err != nil {
		return nil, errors.InternalServer(reason.ApiKeyInvalid).WithError(err).WithStack()
	}

	apiKey := &entities.ApiKey{}
	_ = copier.Copy(apiKey, req)
	apiKey.ID = uid.ID()
	apiKey.Key = hmacedKey
	apiKey.KeyPreview = key[:6]

	err = s.apiKeyRepo.AddApiKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	return &schemas.AddApiKeyResp{
		FriendlyName: apiKey.FriendlyName,
		Key:          fullApiKey,
		KeyPreview:   apiKey.KeyPreview,
		CreatedAt:    apiKey.CreatedAt.Unix(),
	}, nil
}

func (s *ApiKeyService) GetApiKeyWithKey(
	ctx context.Context,
	key string,
) (*schemas.GetApiKeyResp, error) {
	parts := strings.SplitN(key, ".", 3)
	if len(parts) != 3 {
		return nil, errors.Unauthorized(reason.ApiKeyInvalid)
	}

	// Decode the provided signature.
	gotSig := parts[2]
	gotSigBytes, err := converter.FromURLB64(gotSig)
	if err != nil {
		return nil, errors.Unauthorized(reason.ApiKeyInvalid)
	}

	gotKey := parts[0]
	if gotKey == "" {
		return nil, errors.Unauthorized(reason.ApiKeyInvalid)
	}

	gotUserID := parts[1]
	if gotUserID == "" {
		return nil, errors.Unauthorized(reason.ApiKeyInvalid)
	}

	expSig, err := crypto.GetHMAC(gotKey+"."+gotUserID, s.cfg.Key.ApiKeySignatureHMACBytes)
	if err != nil {
		return nil, errors.Unauthorized(reason.ApiKeyInvalid)
	}
	expSigBytes, err := converter.FromURLB64(expSig)
	if err != nil {
		return nil, errors.Unauthorized(reason.ApiKeyInvalid)
	}

	if !hmac.Equal(gotSigBytes, expSigBytes) {
		return nil, errors.Unauthorized(reason.ApiKeyInvalid)
	}

	hmacedKey, err := s.generateApiKeyHMAC(gotKey)
	if err != nil {
		return nil, errors.InternalServer(reason.ApiKeyInvalid).WithError(err).WithStack()
	}

	apiKey, exists, err := s.apiKeyRepo.GetApiKeyByKey(ctx, hmacedKey)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.Unauthorized(reason.ApiKeyInvalid)
	}

	resp := &schemas.GetApiKeyResp{}
	_ = copier.Copy(resp, apiKey)
	resp.CreatedAt = apiKey.CreatedAt.Unix()

	return resp, nil
}

func (s *ApiKeyService) generateApiKey(userID string) (string, error) {
	// Create the "key" parts.
	buf := make([]byte, apiKeyBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to rand: %w", err)
	}
	key := converter.ToURLB64(buf)

	// Add the user ID.
	key = key + "." + userID

	// Create the HMAC of the key and the user.
	sig, err := crypto.GetHMAC(key, s.cfg.Key.ApiKeySignatureHMACBytes)
	if err != nil {
		return "", fmt.Errorf("failed to sign: %w", err)
	}

	// Put it all together.
	key = key + "." + sig

	return key, nil
}

func (s *ApiKeyService) generateApiKeyHMAC(apiKey string) (string, error) {
	sig, err := crypto.GetHMAC(apiKey, s.cfg.Key.ApiKeyDatabaseHMACBytes)
	if err != nil {
		return "", fmt.Errorf("failed to sign: %w", err)
	}

	return sig, nil
}
