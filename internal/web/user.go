package web

import (
	"github.com/htquangg/a-wasm/internal/controllers"

	"github.com/labstack/echo/v4"
)

func bindUserApi(g *echo.Group, c *controllers.Controllers) {
	subGroup := g.Group("/users")

	authGroup := subGroup.Group("/auth")
	authGroup.POST("/sign-up", c.User.SignUp)
	authGroup.POST("/create-session", c.User.CreateSRPSession)
	authGroup.POST("/verify-session", c.User.VerifySRPSession)

	srpGroup := subGroup.Group("/srp")
	srpGroup.GET("/attributes", c.User.GetSRPAttributes)
	srpGroup.POST("/setup", c.User.SetupSRP)
	srpGroup.POST("/complete", c.User.CompleteSRPSetup)
}
