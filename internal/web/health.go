package web

import (
	"github.com/htquangg/a-wasm/internal/controllers"

	"github.com/labstack/echo/v4"
)

type healthApi struct{}

func bindHealthApi(g *echo.Group, c *controllers.Controllers) {
	subGroup := g.Group("/healthz")
	subGroup.GET("/", c.Health.CheckHealth)
}
