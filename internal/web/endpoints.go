package web

import (
	"github.com/htquangg/a-wasm/internal/handlers"
	"github.com/labstack/echo/v4"
)

type endpointsApi struct{}

func bindEndpointsApi(g *echo.Group, h *handlers.Handlers) {
	subGroup := g.Group("/endpoints")

	subGroup.POST("/", h.Endpoint.Add)
}
