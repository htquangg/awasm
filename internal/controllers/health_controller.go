package controllers

import (
	"github.com/labstack/echo/v4"

	"github.com/htquangg/awasm/internal/base/handler"
	"github.com/htquangg/awasm/internal/services/health"
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
