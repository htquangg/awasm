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
	"github.com/htquangg/a-wasm/internal/services/deployment"
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

	userID := middleware.GetUserID(ctx, entities.API_KEY)

	req := &schemas.ServeDeploymentReq{
		DeploymentID: deploymentID,
		UserID:       userID,
		URL:          trimmedEndpointFromURL(ctx.Request().URL),
		Method:       ctx.Request().Method,
		Header:       ctx.Request().Header,
		Body:         b,
	}

	resp, err := c.deploymentService.Serve(ctx.Request().Context(), req)
	if err != nil {
		return handler.HandleResponse(ctx, err, resp)
	}

	return middleware.Serve(ctx, int(resp.StatusCode), resp.Response, resp.Header)
}
