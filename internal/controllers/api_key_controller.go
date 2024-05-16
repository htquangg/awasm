package controllers

import (
	"github.com/labstack/echo/v4"

	"github.com/htquangg/a-wasm/internal/base/handler"
	"github.com/htquangg/a-wasm/internal/base/middleware"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/api_key"
)

type ApiKeyController struct {
	apiKeyService *api_key.ApiKeyService
}

func NewApiKeyController(apiKeyService *api_key.ApiKeyService) *ApiKeyController {
	return &ApiKeyController{
		apiKeyService: apiKeyService,
	}
}

func (c *ApiKeyController) Add(ctx echo.Context) error {
	req := &schemas.AddApiKeyReq{}
	if err, errField := handler.BindAndValidate(ctx, req); err != nil {
		return handler.HandleResponse(ctx, err, errField)
	}

	req.UserID = middleware.GetUserID(ctx, entities.JWT)
	resp, err := c.apiKeyService.AddApiKey(ctx.Request().Context(), req)

	return handler.HandleResponse(ctx, err, resp)
}
