package web

import (
	"github.com/htquangg/a-wasm/internal/services"
	"github.com/labstack/echo/v4"
)

type endpointsApi struct{}

func bindEndpointsApi(g *echo.Group, svc *services.Service) {
	api := &endpointsApi{}

	subGroup := g.Group("/endpoints")
	subGroup.GET("", api.getOne)
	subGroup.GET("/:id", api.getOne)
}

func (api *endpointsApi) getAll(e echo.Context) error {
	return nil
}

func (api *endpointsApi) getOne(c echo.Context) error {
	return nil
}
