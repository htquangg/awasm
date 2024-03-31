package controllers

import (
	"github.com/htquangg/a-wasm/internal/base/handler"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/endpoint"

	"github.com/labstack/echo/v4"
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
	if err, errField := handler.BindAndValidate(ctx, req); err != nil {
		return handler.HandleResponse(ctx, err, errField)
	}

	resp, err := c.endpointService.AddEndpoint(ctx.Request().Context(), req)

	return handler.HandleResponse(ctx, err, resp)
}
