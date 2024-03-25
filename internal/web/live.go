package web

import (
	"github.com/htquangg/a-wasm/internal/controllers"

	"github.com/labstack/echo/v4"
)

func bindLiveApi(g *echo.Group, c *controllers.Controllers) {
	subGroup := g.Group("/live")

	subGroup.POST("/publish", c.Live.Publish)
	subGroup.GET("/:endpointID/*", c.Live.Serve)
}
