package web

import (
	"github.com/htquangg/a-wasm/internal/base/middleware"
	"github.com/htquangg/a-wasm/internal/controllers"

	"github.com/labstack/echo/v4"
)

func bindDeploymentsApi(g *echo.Group, _ *controllers.Controllers, _ *middleware.Middleware) {
	_ = g.Group("/deployments")
}
