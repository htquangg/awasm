package web

import (
	"github.com/labstack/echo/v4"

	"github.com/htquangg/awasm/internal/base/middleware"
	"github.com/htquangg/awasm/internal/controllers"
)

func bindPreviewApi(g *echo.Group, c *controllers.Controllers, mws *middleware.Middleware) {
	subGroup := g.Group("/preview")

	subGroup.Any("/:deploymentID/*", c.Preview.Serve, mws.ApiKey.RequireApiKey)
}
