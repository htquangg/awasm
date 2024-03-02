package controllers

import (
	"io"

	"github.com/htquangg/a-wasm/internal/handler"
	"github.com/htquangg/a-wasm/internal/reason"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/deployment"
	"github.com/segmentfault/pacman/errors"

	"github.com/labstack/echo/v4"
)

type DeploymentController struct {
	deploymentService *deployment.DeploymentService
}

func NewDeploymentController(deploymentService *deployment.DeploymentService) *DeploymentController {
	return &DeploymentController{
		deploymentService: deploymentService,
	}
}

func (c *DeploymentController) Add(ctx echo.Context) error {
	endpointID := ctx.Param("id")

	// TODO: validate the contents of the blob and limit maximum blob size
	b, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return errors.BadRequest(reason.RequestFormatError)
	}

	if len(b) == 0 {
		return errors.BadRequest(reason.RequestFormatError)
	}

	req := &schemas.AddDeploymentReq{
		EndpointID: endpointID,
		Data:       b,
	}

	resp, err := c.deploymentService.Add(ctx.Request().Context(), req)

	return handler.HandleResponse(ctx, err, resp)
}
