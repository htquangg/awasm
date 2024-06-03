package middleware

import (
	"encoding/json"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/cache"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/api_key"
)

const (
	// ApiKeyHeader is the authorization header required for APIKey protected requests.
	ApiKeyHeader = "X-AWASM-API-KEY"
)

type ApiKeyMiddleware struct {
	cfg           *config.Config
	cache         cache.Cacher
	apiKeyService *api_key.ApiKeyService
}

func NewApiKeyMiddleware(
	cfg *config.Config,
	cacher cache.Cacher,
	apiKeyService *api_key.ApiKeyService,
) *ApiKeyMiddleware {
	return &ApiKeyMiddleware{
		cfg:           cfg,
		cache:         cacher,
		apiKeyService: apiKeyService,
	}
}

func (m *ApiKeyMiddleware) RequireApiKey(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		key, err := m.extractApiKey(ctx)
		if err != nil {
			return err
		}

		cacheKey := &cache.Key{
			Namespace: constants.AuthorizedApiKeyNameSpaceCache,
			Key:       key,
		}
		apiKeyBytes, exists, err := m.cache.Fetch(
			ctx.Request().Context(),
			cacheKey,
			constants.AuthorizedApiKeyTTLCache,
			func() (interface{}, error) {
				return m.apiKeyService.GetApiKeyWithKey(ctx.Request().Context(), key)
			},
		)
		if err != nil {
			return err
		}
		if !exists {
			return errors.Unauthorized(reason.ApiKeyInvalid)
		}

		apiKey := &schemas.GetApiKeyResp{}
		err = json.Unmarshal(apiKeyBytes, apiKey)
		if err != nil {
			return errors.Unauthorized(reason.ApiKeyInvalid).WithError(err).WithStack()
		}
		withApiKey(ctx, apiKey)

		return next(ctx)
	}
}

func (m *ApiKeyMiddleware) extractApiKey(ctx echo.Context) (string, error) {
	apiKey := strings.TrimSpace(ctx.Request().Header.Get(ApiKeyHeader))
	if apiKey == "" {
		return "", errors.Unauthorized(reason.ApiKeyRequired)
	}

	return apiKey, nil
}
