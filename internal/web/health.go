package web

import (
	"github.com/labstack/echo/v4"

	"github.com/htquangg/awasm/internal/controllers"
)

func bindHealthApi(g *echo.Group, c *controllers.Controllers) {
	subGroup := g.Group("/healthz")
	subGroup.GET("", c.Health.CheckHealth)
}
