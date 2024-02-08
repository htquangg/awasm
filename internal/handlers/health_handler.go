package handlers

import (
	"github.com/htquangg/a-wasm/internal/services/health"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	healthService *health.HealthService
}

func NewHealthHandler(healthService *health.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

func (h *HealthHandler) CheckHealth(c echo.Context) error {
	resp, err := h.healthService.CheckHealth()
	if err != nil {
		return handleBadRequestResp(c, err, "")
	}

	return handleSuccessResp(c, resp)
}
