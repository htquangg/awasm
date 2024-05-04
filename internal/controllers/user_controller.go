package controllers

import (
	"github.com/labstack/echo/v4"

	"github.com/htquangg/a-wasm/internal/base/handler"
	"github.com/htquangg/a-wasm/internal/base/middleware"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/user"
	"github.com/htquangg/a-wasm/pkg/network"
)

type UserController struct {
	userService *user.UserService
}

func NewUserController(userService *user.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) CheckAuth(ctx echo.Context) error {
	return handler.HandleResponse(ctx, nil, nil)
}

func (c *UserController) GetSRPAttribute(ctx echo.Context) error {
	req := &schemas.GetSRPAttributeReq{}
	if err, errField := handler.BindAndValidate(ctx, req); err != nil {
		return handler.HandleResponse(ctx, err, errField)
	}

	resp, err := c.userService.GetSRPAttribute(ctx.Request().Context(), req)

	return handler.HandleResponse(ctx, err, resp)
}

func (c *UserController) BeginEmailSignupProcess(ctx echo.Context) error {
	req := &schemas.BeginEmailSignupProcessReq{}
	if err, errField := handler.BindAndValidate(ctx, req); err != nil {
		return handler.HandleResponse(ctx, err, errField)
	}

	err := c.userService.BeginEmailSignupProcess(ctx.Request().Context(), req)

	return handler.HandleResponse(ctx, err, nil)
}

func (c *UserController) VerifyEmailSignup(ctx echo.Context) error {
	req := &schemas.VerifyEmailSignupReq{}
	if err, errField := handler.BindAndValidate(ctx, req); err != nil {
		return handler.HandleResponse(ctx, err, errField)
	}

	resp, err := c.userService.VerifyEmailSignup(ctx.Request().Context(), req)

	return handler.HandleResponse(ctx, err, resp)
}

func (c *UserController) SetupSRPAccountSignup(ctx echo.Context) error {
	req := &schemas.SetupSRPAccountSignupReq{}
	if err, errField := handler.BindAndValidate(ctx, req); err != nil {
		return handler.HandleResponse(ctx, err, errField)
	}

	userID := middleware.GetUserID(ctx)
	resp, err := c.userService.SetupSRPAccountSignup(ctx.Request().Context(), userID, req)

	return handler.HandleResponse(ctx, err, resp)
}

func (c *UserController) CompleteEmailAccountSignup(ctx echo.Context) error {
	req := &schemas.CompleteEmailSignupReq{}
	if err, errField := handler.BindAndValidate(ctx, req); err != nil {
		return handler.HandleResponse(ctx, err, errField)
	}
	req.IP = network.GetClientIP(ctx)
	req.UserAgent = ctx.Request().UserAgent()

	resp, err := c.userService.CompleteEmailAccountSignup(ctx.Request().Context(), middleware.GetUser(ctx), req)

	return handler.HandleResponse(ctx, err, resp)
}

func (c *UserController) ChallengeEmailLogin(ctx echo.Context) error {
	req := &schemas.ChallengeEmailLoginReq{}
	if err, errField := handler.BindAndValidate(ctx, req); err != nil {
		return handler.HandleResponse(ctx, err, errField)
	}

	resp, err := c.userService.ChallengeEmailLogin(ctx.Request().Context(), req)

	return handler.HandleResponse(ctx, err, resp)
}

func (c *UserController) VerifyEmailLogin(ctx echo.Context) error {
	req := &schemas.VerifyEmailLoginReq{}
	if err, errField := handler.BindAndValidate(ctx, req); err != nil {
		return handler.HandleResponse(ctx, err, errField)
	}

	resp, err := c.userService.VerifyEmailLogin(ctx.Request().Context(), req)

	return handler.HandleResponse(ctx, err, resp)
}
