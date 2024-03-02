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

func (c *PreviewController) Serve(ctx echo.Context) error {
	deploymentID := ctx.Param("deploymentID")

	b, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return handler.HandleResponse(ctx,
			errors.
				InternalServer(reason.UnknownError).
				WithError(err).
				WithStack(),
			nil)
	}

	req := &schemas.ServeDeploymentReq{
		DeploymentID: deploymentID,
		URL:          trimmedEndpointFromURL(ctx.Request().URL),
		Method:       ctx.Request().Method,
		Header:       ctx.Request().Header,
		Body:         b,
	}

	resp, err := c.deploymentService.Serve(ctx.Request().Context(), req)

	return handler.HandleResponse(ctx, err, resp)
}
