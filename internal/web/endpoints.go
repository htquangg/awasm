package web

import (
	"github.com/htquangg/a-wasm/internal/controllers"
	"github.com/labstack/echo/v4"
)

type endpointsApi struct{}

func bindEndpointsApi(g *echo.Group, h *controllers.Controllers) {
	subGroup := g.Group("/endpoints")

	subGroup.POST("/", h.Endpoint.Add)
}
