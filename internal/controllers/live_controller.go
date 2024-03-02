package controllers

import (
	"io"

	"github.com/htquangg/a-wasm/internal/handler"
	"github.com/htquangg/a-wasm/internal/reason"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/endpoint"
	"github.com/labstack/echo/v4"
	"github.com/segmentfault/pacman/errors"
)

type LiveController struct {
	endpointService *endpoint.EndpointService
}

func NewLiveController(endpointService *endpoint.EndpointService) *LiveController {
	return &LiveController{
		endpointService: endpointService,
	}
}

func (c *LiveController) Publish(ctx echo.Context) error {
	req := &schemas.PublishEndpointReq{}
	if err := ctx.Bind(req); err != nil {
		return errors.BadRequest(reason.RequestFormatError)
	}

	resp, err := c.endpointService.Publish(ctx.Request().Context(), req)

	return handler.HandleResponse(ctx, err, resp)
}

func (c *LiveController) Serve(ctx echo.Context) error {
	endpointID := ctx.Param("endpointID")

	b, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return handler.HandleResponse(ctx,
			errors.
				InternalServer(reason.UnknownError).
				WithError(err).
				WithStack(),
			nil)
	}

	req := &schemas.ServeEndpointReq{
		EndpointID: endpointID,
		URL:        trimmedEndpointFromURL(ctx.Request().URL),
		Method:     ctx.Request().Method,
		Header:     ctx.Request().Header,
		Body:       b,
	}

	resp, err := c.endpointService.Serve(ctx.Request().Context(), req)

	return handler.HandleResponse(ctx, err, resp)
}
