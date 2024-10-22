package middleware

import (
	"regexp"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/awasm/config"
	"github.com/htquangg/awasm/internal/base/reason"
	"github.com/htquangg/awasm/internal/entities"
	"github.com/htquangg/awasm/internal/services/session"
	"github.com/htquangg/awasm/internal/services/user"
)

var bearerRegexp = regexp.MustCompile(`^(?:B|b)earer (\S+$)`)

type AuthMiddleware struct {
	cfg         *config.Config
	userRepo    user.UserRepo
	sessionRepo session.SessionRepo
}

func NewAuthMiddleware(
	cfg *config.Config,
	userRepo user.UserRepo,
	sessionRepo session.SessionRepo,
) *AuthMiddleware {
	return &AuthMiddleware{
		cfg:         cfg,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (m *AuthMiddleware) RequireAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token, err := m.extractBearerToken(ctx)
		if err != nil {
			return err
		}

		err = m.parseJWTClaims(ctx, token)
		if err != nil {
			return err
		}

		err = m.requireAuthentication(ctx)
		if err != nil {
			return err
		}

		return next(ctx)
	}
}

func (m *AuthMiddleware) RequireSignUpAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token, err := m.extractBearerToken(ctx)
		if err != nil {
			return err
		}

		err = m.parseJWTClaims(ctx, token)
		if err != nil {
			return err
		}

		err = m.requireSignUpAuthentication(ctx)
		if err != nil {
			return err
		}

		return next(ctx)
	}
}

func (m *AuthMiddleware) RequireSignupToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		claims := getClaims(ctx)
		if claims == nil {
			return errors.Unauthorized(reason.InvalidTokenError)
		}

		if claims.CommonTokenClaims.GetScope() != entities.SignupTokenScope {
			return errors.Forbidden(reason.InvalidScopeError)
		}

		return next(ctx)
	}
}

func (m *AuthMiddleware) extractBearerToken(ctx echo.Context) (string, error) {
	token := ctx.Request().Header.Get("Authorization")
	if token == "" {
		token = ctx.Request().URL.Query().Get("Authorization")
	}

	matches := bearerRegexp.FindStringSubmatch(token)
	if len(matches) != 2 {
		return "", errors.Unauthorized(reason.InvalidTokenError)
	}

	return matches[1], nil
}

func (m *AuthMiddleware) parseJWTClaims(echoCtx echo.Context, bearer string) error {
	p := jwt.Parser{
		ValidMethods: []string{jwt.SigningMethodHS256.Name},
	}
	token, err := p.ParseWithClaims(
		bearer,
		&entities.AccessTokenClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return (m.cfg.JWT.SecretBytes), nil
		},
	)
	if err != nil {
		return errors.Unauthorized(reason.InvalidTokenError)
	}

	withToken(echoCtx, token)

	return nil
}

func (m *AuthMiddleware) requireAuthentication(ctx echo.Context) error {
	claims := getClaims(ctx)

	if claims == nil {
		return errors.Unauthorized(reason.InvalidTokenError)
	}

	if claims.Subject == "" {
		return errors.Unauthorized(reason.InvalidTokenError)
	}

	userID := claims.Subject
	user, exists, err := m.userRepo.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.Forbidden(reason.UserNotFound)
	}
	withUser(ctx, user)

	if claims.SessionID == "" {
		return errors.Unauthorized(reason.InvalidTokenError)
	}

	session, exists, err := m.sessionRepo.GetSessionByID(
		ctx.Request().Context(),
		claims.SessionID,
		false,
	)
	if err != nil {
		return err
	}
	if !exists {
		return errors.Forbidden(reason.SessionNotFound)
	}
	withSession(ctx, session)

	return nil
}

func (m *AuthMiddleware) requireSignUpAuthentication(ctx echo.Context) error {
	claims := getClaims(ctx)

	if claims == nil {
		return errors.Unauthorized(reason.InvalidTokenError)
	}

	if claims.Subject == "" {
		return errors.Unauthorized(reason.InvalidTokenError)
	}

	userID := claims.Subject
	user, exists, err := m.userRepo.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.Forbidden(reason.UserNotFound)
	}
	withUser(ctx, user)

	return nil
}
