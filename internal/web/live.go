package web

import (
	"github.com/htquangg/a-wasm/internal/controllers"
	"github.com/labstack/echo/v4"
)

func bindLiveApi(g *echo.Group, h *controllers.Controllers) {
	subGroup := g.Group("/live")

	subGroup.POST("/publish", h.Live.Publish)
	subGroup.GET("/:endpointID/*", h.Live.Serve)
}
