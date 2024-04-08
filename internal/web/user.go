package web

import (
	"github.com/htquangg/a-wasm/internal/base/middleware"
	"github.com/htquangg/a-wasm/internal/controllers"

	"github.com/labstack/echo/v4"
)

func bindUserApi(g *echo.Group, c *controllers.Controllers, mws *middleware.Middleware) {
	subGroup := g.Group("/users")

	authGroup := subGroup.Group("/auth")

	publicAuthEmailGroup := authGroup.Group("/email")
	publicAuthEmailGroup.POST("/signup", c.User.BeginEmailSignupProcess)
	publicAuthEmailGroup.POST("/verify", c.User.VerifyEmailSignup)
	publicAuthEmailGroup.POST("/create-session", c.User.CreateSRPSession, mws.Auth.RequireAuthentication)
	publicAuthEmailGroup.POST("/verify-session", c.User.VerifySRPSession)

	privateAuthEmailGroup := authGroup.Group("/email", mws.Auth.RequireAuthentication)
	privateAuthEmailGroup.POST("/setup-srp", c.User.SetupSRPAccountSignup, mws.Auth.RequireSignupToken)
	privateAuthEmailGroup.POST("/complete", c.User.CompleteEmailAccountSignup, mws.Auth.RequireSignupToken)
}
