package web

import (
	"github.com/htquangg/a-wasm/internal/controllers"

	"github.com/labstack/echo/v4"
)

func bindPreviewApi(g *echo.Group, c *controllers.Controllers) {
	subGroup := g.Group("/preview")

	subGroup.Any("/:deploymentID/*", c.Preview.Serve)
}
