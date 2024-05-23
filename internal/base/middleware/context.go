package middleware

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/schemas"
)

const (
	tokenKey   = "ctx-req-jwt"
	userKey    = "ctx-req-user"
	sessionKey = "ctx-req-session"
	apiKey     = "ctx-api-key"
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

func GetUserID(ctx echo.Context, authMode entities.AuthMode) string {
	userID := ""

	switch authMode {
	case entities.AuthModeJwt:
		user := getUser(ctx)
		if user != nil {
			userID = user.ID
		}

	case entities.AuthModeApiKey:
		apiKey := getApiKey(ctx)
		if apiKey != nil {
			userID = apiKey.UserID
		}

	default:
		userID = ""
	}

	return userID
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

func withApiKey(ctx echo.Context, k *schemas.GetApiKeyResp) {
	ctx.Set(apiKey, k)
}

func GetApiKey(ctx echo.Context) *schemas.GetApiKeyResp {
	return getApiKey(ctx)
}

func getApiKey(ctx echo.Context) *schemas.GetApiKeyResp {
	if ctx == nil {
		return nil
	}
	obj := ctx.Get(apiKey)
	if obj == nil {
		return nil
	}
	return obj.(*schemas.GetApiKeyResp)
}
