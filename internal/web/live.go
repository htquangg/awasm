package web

import (
	"github.com/labstack/echo/v4"

	"github.com/htquangg/a-wasm/internal/base/middleware"
	"github.com/htquangg/a-wasm/internal/controllers"
)

func bindLiveApi(g *echo.Group, c *controllers.Controllers, mws *middleware.Middleware) {
	subGroup := g.Group("/live")

	subGroup.POST("/publish", c.Live.Publish, mws.Auth.RequireAuthentication)
	subGroup.Any("/:endpointID/*", c.Live.Serve)
}
