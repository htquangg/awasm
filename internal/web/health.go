package web

import (
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/labstack/echo/v4"
)

type healthApi struct{}

func bindHealthApi(g *echo.Group) {
	api := &healthApi{}

	subGroup := g.Group("/healthz")
	subGroup.GET("/", api.checkHealth)
}

func (api *healthApi) checkHealth(c echo.Context) error {
	return handleSuccessResp(c, &schemas.CheckHealthResp{Msg: "API is live!!!"})
}
