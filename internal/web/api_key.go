package web

import (
	"github.com/labstack/echo/v4"

	"github.com/htquangg/a-wasm/internal/base/middleware"
	"github.com/htquangg/a-wasm/internal/controllers"
)

func bindApiKeyApi(g *echo.Group, c *controllers.Controllers, mws *middleware.Middleware) {
	subGroup := g.Group("/api-keys")

	subGroup.POST("", c.ApiKey.Add, mws.Auth.RequireAuthentication)
}
