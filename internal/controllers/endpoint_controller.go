package controllers

import (
	"github.com/labstack/echo/v4"

	"github.com/htquangg/awasm/internal/base/handler"
	"github.com/htquangg/awasm/internal/base/middleware"
	"github.com/htquangg/awasm/internal/entities"
	"github.com/htquangg/awasm/internal/schemas"
	"github.com/htquangg/awasm/internal/services/endpoint"
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

	req.UserID = middleware.GetUserID(ctx, entities.AuthModeJwt)
	resp, err := c.endpointService.AddEndpoint(ctx.Request().Context(), req)

	return handler.HandleResponse(ctx, err, resp)
}
