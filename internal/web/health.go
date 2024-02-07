package web

import (
	"github.com/htquangg/a-wasm/internal/services/health"

	"github.com/labstack/echo/v4"
)

type healthApi struct{}

func bindHealthApi(g *echo.Group) {
	api := &healthApi{}

	subGroup := g.Group("/healthz")
	subGroup.GET("/", api.checkHealth)
}

func (api *healthApi) checkHealth(c echo.Context) error {
	return handleSuccessResp(c, &health.CheckHealthResp{Msg: "API is live!!!"})
}
