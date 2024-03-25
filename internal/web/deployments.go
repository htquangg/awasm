package web

import (
	"github.com/htquangg/a-wasm/internal/controllers"

	"github.com/labstack/echo/v4"
)

func bindDeploymentsApi(g *echo.Group, c *controllers.Controllers) {
	_ = g.Group("/deployments")
}
