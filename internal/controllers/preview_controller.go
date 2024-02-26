package controllers

import (
	"io"

	"github.com/htquangg/a-wasm/internal/handler"
	"github.com/htquangg/a-wasm/internal/reason"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/deployment"

	"github.com/labstack/echo/v4"
	"github.com/segmentfault/pacman/errors"
)

type PreviewController struct {
	deploymentService *deployment.DeploymentService
}

func NewPreviewController(deploymentService *deployment.DeploymentService) *PreviewController {
	return &PreviewController{
		deploymentService: deploymentService,
	}
}

func (h *PreviewController) Serve(c echo.Context) error {
	deploymentID := c.Param("deploymentID")

	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return handler.HandleResponse(c,
			errors.
				InternalServer(reason.UnknownError).
				WithError(err).
				WithStack(),
			nil)
	}

	req := &schemas.ServePreviewReq{
		DeploymentID: deploymentID,
		URL:          trimmedEndpointFromURL(c.Request().URL),
		Method:       c.Request().Method,
		Header:       c.Request().Header,
		Body:         b,
	}

	resp, err := h.deploymentService.Serve(c.Request().Context(), req)

	return handler.HandleResponse(c, err, resp)
}
