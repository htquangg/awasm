package controllers

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/awasm/internal/base/handler"
	"github.com/htquangg/awasm/internal/base/middleware"
	"github.com/htquangg/awasm/internal/base/reason"
	"github.com/htquangg/awasm/internal/entities"
	"github.com/htquangg/awasm/internal/schemas"
	"github.com/htquangg/awasm/internal/services/deployment"
)

type DeploymentController struct {
	deploymentService *deployment.DeploymentService
}

func NewDeploymentController(
	deploymentService *deployment.DeploymentService,
) *DeploymentController {
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
		UserID:     middleware.GetUserID(ctx, entities.AuthModeJwt),
		Data:       b,
	}

	resp, err := c.deploymentService.AddDeployment(ctx.Request().Context(), req)

	return handler.HandleResponse(ctx, err, resp)
}
