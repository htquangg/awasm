package session

import (
	"context"
	"time"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/schemas"

	"github.com/golang-jwt/jwt"
	"github.com/segmentfault/pacman/errors"
)

type (
	SessionRepo interface {
		CreateRefreshToken(
			ctx context.Context,
			userID string,
			authenticationMethod entities.AuthenticationMethod,
			params *entities.GrantParams,
		) (*entities.RefreshToken, error)
		GetSessionByID(ctx context.Context, id string) (*entities.Session, bool, error)
	}

	SessionService struct {
		cfg         *config.Config
		sessionRepo SessionRepo
	}
)

func NewSessionService(cfg *config.Config, sessionRepo SessionRepo) *SessionService {
	return &SessionService{
		cfg:         cfg,
		sessionRepo: sessionRepo,
	}
}

func (s *SessionService) IssueRefreshToken(
	ctx context.Context,
	user *entities.User,
	authenticationMethod entities.AuthenticationMethod,
	params *entities.GrantParams,
) (*schemas.AccessTokenResp, error) {
	refreshToken, err := s.sessionRepo.CreateRefreshToken(ctx, user.ID, authenticationMethod, params)
	if err != nil {
		return nil, err
	}

	accessToken, expiresAt, err := s.generateAccessToken(ctx, user, refreshToken.SessionID)
	if err != nil {
		return nil, err
	}

	return &schemas.AccessTokenResp{
		CommonTokenResp: schemas.CommonTokenResp{
			AccessToken: accessToken,
			TokenType:   "bearer",
			ExpiresIn:   s.cfg.JWT.Exp,
			ExpiresAt:   expiresAt,
		},
		RefreshToken: refreshToken.Token,
	}, nil
}

func (s *SessionService) IssueSignupToken(
	ctx context.Context,
	user *entities.User,
	params *entities.GrantParams,
) (*schemas.CommonTokenResp, error) {
	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(constants.ExpiresInTokenSignup).Unix()

	claims := &entities.CommonTokenClaims{
		StandardClaims: jwt.StandardClaims{Subject: user.ID, IssuedAt: issuedAt.Unix(), ExpiresAt: expiresAt},
		Email:          user.Email,
		Scope:          entities.SignupTokenScope.Ptr(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(s.cfg.JWT.SecretBytes)
	if err != nil {
		return nil, err
	}

	return &schemas.CommonTokenResp{
		AccessToken: signed,
		TokenType:   "bearer",
		ExpiresIn:   int(constants.ExpiresInTokenSignup.Seconds()),
		ExpiresAt:   expiresAt,
	}, nil
}

func (s *SessionService) generateAccessToken(
	ctx context.Context,
	user *entities.User,
	sessionID string,
) (string, int64, error) {
	if sessionID == "" {
		return "", 0, errors.InternalServer(reason.RequiredSession)
	}

	session, exists, err := s.sessionRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return "", 0, err
	}
	if !exists {
		return "", 0, errors.BadRequest(reason.SessionNotFound)
	}

	aal, amr, err := session.CalculateAALAndAMR(user)
	if err != nil {
		return "", 0, err
	}

	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(time.Second * time.Duration(s.cfg.JWT.Exp)).Unix()

	claims := &entities.AccessTokenClaims{
		CommonTokenClaims: entities.CommonTokenClaims{
			StandardClaims: jwt.StandardClaims{Subject: user.ID, IssuedAt: issuedAt.Unix(), ExpiresAt: expiresAt},
			Email:          user.Email,
			Scope:          entities.AccessTokenScope.Ptr(),
		},
		AuthenticatorAssuranceLevel:   aal.String(),
		AuthenticationMethodReference: amr,
		SessionID:                     sessionID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(s.cfg.JWT.SecretBytes)
	if err != nil {
		return "", 0, err
	}

	return signed, expiresAt, nil
}
