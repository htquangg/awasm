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
	signupEmailGroup.POST("/challenge", c.User.BeginEmailSignupProcess)
	signupEmailGroup.POST("/verify", c.User.VerifyEmailSignup)
	signupEmailGroup.POST(
		"/complete",
		c.User.CompleteEmailAccountSignup,
		mws.Auth.RequireAuthentication,
		mws.Auth.RequireSignupToken,
	)

	loginEmailGroup := emailAuthGroup.Group("/login")
	loginEmailGroup.POST("/challenge", c.User.ChallengeEmailLogin)
	loginEmailGroup.POST("/verify", c.User.VerifyEmailLogin)

	// srp group
	srpGroup := subGroup.Group("/srp")
	srpGroup.GET("/attributes", c.User.GetSRPAttribute)
	srpGroup.POST("/setup", c.User.SetupSRPAccountSignup, mws.Auth.RequireAuthentication, mws.Auth.RequireSignupToken)
}
