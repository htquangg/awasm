package middleware

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/htquangg/a-wasm/internal/entities"
)

const (
	tokenKey   = "ctx-req-jwt"
	userKey    = "ctx-req-user"
	sessionKey = "ctx-req-session"
)

func withToken(ctx echo.Context, token *jwt.Token) {
	ctx.Set(tokenKey, token)
}

func GetToken(ctx echo.Context) *jwt.Token {
	return getToken(ctx)
}

func getToken(ctx echo.Context) *jwt.Token {
	obj := ctx.Get(tokenKey)
	if obj == nil {
		return nil
	}

	return obj.(*jwt.Token)
}

func getClaims(ctx echo.Context) *entities.AccessTokenClaims {
	token := getToken(ctx)
	if token == nil {
		return nil
	}
	return token.Claims.(*entities.AccessTokenClaims)
}

func withUser(ctx echo.Context, user *entities.User) {
	ctx.Set(userKey, user)
}

func GetUser(ctx echo.Context) *entities.User {
	return getUser(ctx)
}

func GetUserID(ctx echo.Context) string {
	user := getUser(ctx)
	if user == nil {
		return ""
	}

	return user.ID
}

func getUser(ctx echo.Context) *entities.User {
	if ctx == nil {
		return nil
	}
	obj := ctx.Get(userKey)
	if obj == nil {
		return nil
	}
	return obj.(*entities.User)
}

func withSession(ctx echo.Context, session *entities.Session) {
	ctx.Set(sessionKey, session)
}

func GetSession(ctx echo.Context) *entities.Session {
	return getSession(ctx)
}

func getSession(ctx echo.Context) *entities.Session {
	if ctx == nil {
		return nil
	}
	obj := ctx.Get(sessionKey)
	if obj == nil {
		return nil
	}
	return obj.(*entities.Session)
}
