package controllers

import (
	"github.com/htquangg/a-wasm/internal/base/handler"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/user"
	"github.com/htquangg/a-wasm/pkg/network"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService *user.UserService
}

func NewUserController(userService *user.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) SignUp(ctx echo.Context) error {
	req := &schemas.SignUpReq{}
	if err, errField := handler.BindAndValidate(ctx, req); err != nil {
		return handler.HandleResponse(ctx, err, errField)
	}
	req.IP = network.GetClientIP(ctx)
	req.UserAgent = ctx.Request().UserAgent()

	resp, err := c.userService.SignUp(ctx.Request().Context(), req)

	return handler.HandleResponse(ctx, err, resp)
}

func (c *UserController) CreateSRPSession(ctx echo.Context) error {
	return nil
}

func (c *UserController) VerifySRPSession(ctx echo.Context) error {
	return nil
}

func (c *UserController) GetSRPAttributes(ctx echo.Context) error {
	return nil
}

func (c *UserController) SetupSRP(ctx echo.Context) error {
	return nil
}

func (c *UserController) CompleteSRPSetup(ctx echo.Context) error {
	return nil
}
