package web

import (
	"github.com/htquangg/a-wasm/internal/base/middleware"
	"github.com/htquangg/a-wasm/internal/controllers"

	"github.com/labstack/echo/v4"
)

func bindUserApi(g *echo.Group, c *controllers.Controllers, mws *middleware.Middleware) {
	subGroup := g.Group("/users")

	authGroup := subGroup.Group("/auth")

	publicAuthGroup := authGroup.Group("")
	publicAuthGroup.POST("/verify-email", c.User.VerifyEmail)
	publicAuthGroup.POST("/create-session", c.User.CreateSRPSession)
	publicAuthGroup.POST("/verify-session", c.User.VerifySRPSession)

	srpGroup := subGroup.Group("/srp")

	publicSrpGroup := srpGroup.Group("")
	publicSrpGroup.GET("/attributes", c.User.GetSRPAttributes)

	privateSrpGroup := srpGroup.Group("", mws.Auth.RequireAuthentication)
	privateSrpGroup.POST("/setup", c.User.SetupSRP)
	privateSrpGroup.POST("/complete", c.User.CompleteSRPSetup)
}
