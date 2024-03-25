package web

import (
	"github.com/htquangg/a-wasm/internal/controllers"

	"github.com/labstack/echo/v4"
)

func bindEndpointsApi(g *echo.Group, c *controllers.Controllers) {
	subGroup := g.Group("/endpoints")

	subGroup.POST("/", c.Endpoint.Add)
	subGroup.POST("/:id/deployments", c.Deployment.Add)
}
