package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/services/api_key"
)

const (
	// APIKeyHeader is the authorization header required for APIKey protected requests.
	APIKeyHeader = "X-AWASM-API-Key"
)

type ApiKeyMiddleware struct {
	cfg           *config.Config
	apiKeyService *api_key.ApiKeyService
}

func NewApiKeyMiddleware(cfg *config.Config, apiKeyService *api_key.ApiKeyService) *ApiKeyMiddleware {
	return &ApiKeyMiddleware{
		cfg:           cfg,
		apiKeyService: apiKeyService,
	}
}

func (m *ApiKeyMiddleware) RequireApiKey(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		key, err := m.extractApiKey(ctx)
		if err != nil {
			return err
		}

		_, err = m.apiKeyService.GetApiKeyWithKey(ctx.Request().Context(), key)
		if err != nil {
			return err
		}

		return next(ctx)
	}
}

func (m *ApiKeyMiddleware) extractApiKey(ctx echo.Context) (string, error) {
	apiKey := strings.TrimSpace(ctx.Request().Header.Get(APIKeyHeader))
	if apiKey == "" {
		return "", errors.Unauthorized(reason.ApiKeyRequired)
	}

	return apiKey, nil
}
