package middleware

import (
	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/repos"
)

type Middleware struct {
	cfg  *config.Config
	Auth *AuthMiddleware
}

func NewMiddleware(cfg *config.Config, repos *repos.Repos) *Middleware {
	authMiddleware := NewAuthMiddleware(cfg, repos.User, repos.Session)
	return &Middleware{
		cfg:  cfg,
		Auth: authMiddleware,
	}
}
