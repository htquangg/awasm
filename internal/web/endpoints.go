package web

import (
	"github.com/htquangg/a-wasm/internal/base/middleware"
	"github.com/htquangg/a-wasm/internal/controllers"

	"github.com/labstack/echo/v4"
)

func bindEndpointsApi(g *echo.Group, c *controllers.Controllers, mws *middleware.Middleware) {
	subGroup := g.Group("/endpoints")

	subGroup.POST("", c.Endpoint.Add, mws.Auth.RequireAuthentication)
	subGroup.POST("/:id/deployments", c.Deployment.Add, mws.Auth.RequireAuthentication)
}
