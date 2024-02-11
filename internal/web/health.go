package web

import (
	"github.com/htquangg/a-wasm/internal/controllers"

	"github.com/labstack/echo/v4"
)

type healthApi struct{}

func bindHealthApi(g *echo.Group, h *controllers.Controllers) {
	subGroup := g.Group("/healthz")
	subGroup.GET("/", h.Health.CheckHealth)
}
