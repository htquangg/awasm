package web

import (
	"github.com/labstack/echo/v4"

	"github.com/htquangg/awasm/internal/base/middleware"
	"github.com/htquangg/awasm/internal/controllers"
)

func bindEndpointsApi(g *echo.Group, c *controllers.Controllers, mws *middleware.Middleware) {
	subGroup := g.Group("/endpoints")

	subGroup.POST("", c.Endpoint.Add, mws.Auth.RequireAuthentication)
	subGroup.POST("/:id/deployments", c.Deployment.Add, mws.Auth.RequireAuthentication)
}
