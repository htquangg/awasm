package web

import (
	"github.com/htquangg/a-wasm/internal/base/middleware"
	"github.com/htquangg/a-wasm/internal/controllers"

	"github.com/labstack/echo/v4"
)

func bindUserApi(g *echo.Group, c *controllers.Controllers, mws *middleware.Middleware) {
	subGroup := g.Group("/users")

	// auth group
	authGroup := subGroup.Group("/auth")

	authGroup.POST("/check", c.User.CheckAuth, mws.Auth.RequireAuthentication)

	emailAuthGroup := authGroup.Group("/email")
	signupEmailGroup := emailAuthGroup.Group("/signup")

	publicSignupEmailGroup := signupEmailGroup.Group("")
	publicSignupEmailGroup.POST("/challenge", c.User.BeginEmailSignupProcess)
	publicSignupEmailGroup.POST("/verify", c.User.VerifyEmailSignup)

	privateAuthEmailGroup := signupEmailGroup.Group("", mws.Auth.RequireAuthentication)
	privateAuthEmailGroup.POST("/complete", c.User.CompleteEmailAccountSignup, mws.Auth.RequireSignupToken)

	loginEmailGroup := emailAuthGroup.Group("/login")
	loginEmailGroup.POST("/challenge", c.User.ChallengeEmailLogin)
	loginEmailGroup.POST("/verify", c.User.VerifyEmailLogin)

	// srp group
	srpGroup := subGroup.Group("/srp")

	publicSrpGroup := srpGroup.Group("")
	publicSrpGroup.GET("/attributes", c.User.GetSRPAttribute)

	privateSrpGroup := srpGroup.Group("", mws.Auth.RequireAuthentication)
	privateSrpGroup.POST("/setup", c.User.SetupSRPAccountSignup, mws.Auth.RequireSignupToken)
}
