package network

import "github.com/labstack/echo/v4"

func GetClientIP(ctx echo.Context) string {
	ip := ctx.Request().Header.Get("CF-Connecting-IP")
	if ip == "" {
		ip = ctx.RealIP()
	}
	return ip
}
