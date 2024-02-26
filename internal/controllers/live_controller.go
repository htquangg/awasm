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

func NewLiveController() *LiveController {
	return &LiveController{}
}

func (h *LiveController) Serve(c echo.Context) error {
	endpointID := c.Param("endpointID")

	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return handler.HandleResponse(c,
			errors.
				InternalServer(reason.UnknownError).
				WithError(err).
				WithStack(),
			nil)
	}

	req := &schemas.ServeLiveReq{
		EndpointID: endpointID,
		URL:        trimmedEndpointFromURL(c.Request().URL),
		Method:     c.Request().Method,
		Header:     c.Request().Header,
		Body:       b,
	}

	resp, err := h.endpointService.Serve(c.Request().Context(), req)

	return handler.HandleResponse(c, err, resp)
}
