package web

import (
	"github.com/htquangg/a-wasm/internal/base/middleware"
	"github.com/htquangg/a-wasm/internal/controllers"

	"github.com/labstack/echo/v4"
)

func bindLiveApi(g *echo.Group, c *controllers.Controllers, mws *middleware.Middleware) {
	subGroup := g.Group("/live")

	privateGroup := subGroup.Group("", mws.Auth.RequireAuthentication)
	privateGroup.POST("/publish", c.Live.Publish)

	publicGroup := subGroup.Group("")
	publicGroup.Any("/:endpointID/*", c.Live.Serve)
}
