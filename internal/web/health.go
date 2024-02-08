package web

import (
	"github.com/htquangg/a-wasm/internal/handlers"

	"github.com/labstack/echo/v4"
)

type healthApi struct{}

func bindHealthApi(g *echo.Group, h *handlers.Handlers) {
	subGroup := g.Group("/healthz")
	subGroup.GET("/", h.Health.CheckHealth)
}
