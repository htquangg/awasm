package controllers

import (
	"github.com/htquangg/a-wasm/internal/base/handler"
	"github.com/htquangg/a-wasm/internal/services/health"

	"github.com/labstack/echo/v4"
)

type HealthController struct {
	healthService *health.HealthService
}

func NewHealthController(healthService *health.HealthService) *HealthController {
	return &HealthController{
		healthService: healthService,
	}
}

func (c *HealthController) CheckHealth(ctx echo.Context) error {
	result, err := c.healthService.CheckHealth(ctx.Request().Context())
	if err != nil {
		return handler.HandleResponse(ctx, err, nil)
	}

	return handler.HandleResponse(ctx, nil, result)
}
