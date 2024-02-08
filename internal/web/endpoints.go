package web

import (
	"github.com/labstack/echo/v4"
)

type endpointsApi struct{}

func bindEndpointsApi(g *echo.Group) {
	_ = g.Group("/endpoints")
}
