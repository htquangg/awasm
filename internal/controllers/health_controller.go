package controllers

import (
	"github.com/htquangg/a-wasm/internal/handler"
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

func (h *HealthController) CheckHealth(c echo.Context) error {
	result, err := h.healthService.CheckHealth()
	if err != nil {
		return handler.HandleResponse(c, err, nil)
	}

	return handler.HandleResponse(c, nil, result)
}
