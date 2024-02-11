package handlers

import (
	"github.com/htquangg/a-wasm/internal/handlers/resp"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/endpoint"
	"github.com/htquangg/a-wasm/pkg/uid"

	"github.com/labstack/echo/v4"
	"github.com/segmentfault/pacman/errors"
)

type EndpointHandler struct {
	endpointService *endpoint.EndpointService
}

func NewEndpointHandler(endpointService *endpoint.EndpointService) *EndpointHandler {
	return &EndpointHandler{
		endpointService: endpointService,
	}
}

func (h *EndpointHandler) Add(c echo.Context) error {
	req := &schemas.AddEndpointReq{}
	if err := c.Bind(req); err != nil {
		return errors.BadRequest(resp.RequestFormatError)
	}

	req.ID = uid.ID()

	result, err := h.endpointService.Add(c.Request().Context(), req)

	return resp.HandleResponse(c, err, result)
}
