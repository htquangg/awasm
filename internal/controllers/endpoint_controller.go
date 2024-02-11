package controllers

import (
	"github.com/htquangg/a-wasm/internal/handler"
	"github.com/htquangg/a-wasm/internal/reason"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/endpoint"
	"github.com/htquangg/a-wasm/pkg/uid"

	"github.com/labstack/echo/v4"
	"github.com/segmentfault/pacman/errors"
)

type EndpointController struct {
	endpointService *endpoint.EndpointService
}

func NewEndpointController(endpointService *endpoint.EndpointService) *EndpointController {
	return &EndpointController{
		endpointService: endpointService,
	}
}

func (h *EndpointController) Add(c echo.Context) error {
	req := &schemas.AddEndpointReq{}
	if err := c.Bind(req); err != nil {
		return errors.BadRequest(reason.RequestFormatError)
	}

	req.ID = uid.ID()

	result, err := h.endpointService.Add(c.Request().Context(), req)

	return handler.HandleResponse(c, err, result)
}
