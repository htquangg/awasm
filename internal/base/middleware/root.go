package middleware

import (
	"github.com/htquangg/awasm/config"
	"github.com/htquangg/awasm/internal/base/cache"
	"github.com/htquangg/awasm/internal/repos"
	"github.com/htquangg/awasm/internal/services"
)

type Middleware struct {
	cfg    *config.Config
	Auth   *AuthMiddleware
	ApiKey *ApiKeyMiddleware
}

func NewMiddleware(
	cfg *config.Config,
	cacher cache.Cacher,
	services *services.Sevices,
	repos *repos.Repos,
) *Middleware {
	authMiddleware := NewAuthMiddleware(cfg, repos.User, repos.Session)
	apiKeyMiddleware := NewApiKeyMiddleware(cfg, cacher, services.ApiKey)
	return &Middleware{
		cfg:    cfg,
		Auth:   authMiddleware,
		ApiKey: apiKeyMiddleware,
	}
}
