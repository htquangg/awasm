package middleware

import (
	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/repos"
	"github.com/htquangg/a-wasm/internal/services"
)

type Middleware struct {
	cfg    *config.Config
	Auth   *AuthMiddleware
	ApiKey *ApiKeyMiddleware
}

func NewMiddleware(cfg *config.Config, services *services.Sevices, repos *repos.Repos) *Middleware {
	authMiddleware := NewAuthMiddleware(cfg, repos.User, repos.Session)
	apiKeyMiddleware := NewApiKeyMiddleware(cfg, services.ApiKey)
	return &Middleware{
		cfg:    cfg,
		Auth:   authMiddleware,
		ApiKey: apiKeyMiddleware,
	}
}
