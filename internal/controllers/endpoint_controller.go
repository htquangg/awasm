package controllers

import (
	"github.com/htquangg/a-wasm/internal/base/handler"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/endpoint"

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

func (c *EndpointController) Add(ctx echo.Context) error {
	req := &schemas.AddEndpointReq{}
	if err := ctx.Bind(req); err != nil {
		return errors.BadRequest(reason.RequestFormatError)
	}

	resp, err := c.endpointService.Add(ctx.Request().Context(), req)

	return handler.HandleResponse(ctx, err, resp)
}
