package controllers

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/a-wasm/internal/base/handler"
	"github.com/htquangg/a-wasm/internal/base/middleware"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/endpoint"
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
	if err, errField := handler.BindAndValidate(ctx, req); err != nil {
		return handler.HandleResponse(ctx, err, errField)
	}

	req.UserID = middleware.GetUserID(ctx, entities.AuthModeJwt)
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
	if err != nil {
		return handler.HandleResponse(ctx, err, resp)
	}

	return middleware.Serve(ctx, int(resp.StatusCode), resp.Response, resp.Header)
}
