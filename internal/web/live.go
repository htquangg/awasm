package web

import (
	"github.com/htquangg/a-wasm/internal/controllers"
	"github.com/labstack/echo/v4"
)

func bindLiveApi(g *echo.Group, h *controllers.Controllers) {
	subGroup := g.Group("/live")

	subGroup.GET("/:endpointID", h.Live.Serve)
}
