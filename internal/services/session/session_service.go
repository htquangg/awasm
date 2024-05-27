package session

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/user_common"
	"github.com/htquangg/a-wasm/pkg/crypto"
	"github.com/htquangg/a-wasm/pkg/uid"
)

type (
	SessionRepo interface {
		WithTx(
			parentCtx context.Context,
			f func(ctx context.Context) (interface{}, error),
		) (interface{}, error)
		AddRefreshToken(ctx context.Context, refreshToken *entities.RefreshToken) error
		AddSession(ctx context.Context, session *entities.Session) error
		AddClaimToSession(
			ctx context.Context,
			sessionID string,
			authenticationMethod entities.AuthenticationMethod,
		) error
		UpdateRefreshTokenCols(
			ctx context.Context,
			refreshToken *entities.RefreshToken,
			cols ...string,
		) error
		GetSessionByID(
			ctx context.Context,
			id string,
			forUpdate bool,
		) (*entities.Session, bool, error)
		GetRefreshTokenByToken(
			ctx context.Context,
			token string,
			forUpdate bool,
		) (*entities.RefreshToken, *entities.Session, error)
		GetCurrentlyActiveRefreshTokenBySessionID(
			ctx context.Context,
			sessionID string,
		) (*entities.RefreshToken, bool, error)
	}

	SessionService struct {
		cfg            *config.Config
		sessionRepo    SessionRepo
		userCommonRepo user_common.UserCommonRepo
	}
)

func NewSessionService(
	cfg *config.Config,
	sessionRepo SessionRepo,
	userCommonRepo user_common.UserCommonRepo,
) *SessionService {
	return &SessionService{
		cfg:            cfg,
		sessionRepo:    sessionRepo,
		userCommonRepo: userCommonRepo,
	}
}

func (s *SessionService) IssueRefreshToken(
	ctx context.Context,
	user *entities.User,
	authenticationMethod entities.AuthenticationMethod,
	params *entities.GrantParams,
) (*schemas.AccessTokenResp, error) {
	var refreshToken *entities.RefreshToken
	var accessToken string
	var expiresAt int64

	_, err := s.sessionRepo.WithTx(ctx, func(ctx context.Context) (interface{}, error) {
		var terr error

		refreshToken, terr = s.grantAuthenticatedUser(ctx, user, params)
		if terr != nil {
			return nil, terr
		}

		if terr := s.sessionRepo.AddClaimToSession(ctx, refreshToken.SessionID, authenticationMethod); terr != nil {
			return nil, terr
		}

		accessToken, expiresAt, terr = s.generateAccessToken(ctx, user, refreshToken.SessionID)
		if terr != nil {
			return nil, terr
		}

		return refreshToken, nil
	})
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

func (s *SessionService) GetUserWithRefreshToken(
	ctx context.Context,
	token string,
	forUpdate bool,
) (*entities.User, *entities.RefreshToken, *entities.Session, error) {
	refreshToken, session, err := s.sessionRepo.GetRefreshTokenByToken(ctx, token, forUpdate)
	if err != nil {
		return nil, nil, nil, err
	}

	user, exists, err := s.userCommonRepo.GetUserByID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, nil, nil, err
	}
	if !exists {
		return nil, nil, nil, errors.BadRequest(reason.UserNotFound)
	}

	return user, refreshToken, session, nil
}

func (s *SessionService) GrantRefreshTokenSwap(
	ctx context.Context,
	user *entities.User,
	oldRefreshToken *entities.RefreshToken,
) (*entities.RefreshToken, error) {
	refreshToken := &entities.RefreshToken{}

	_, err := s.sessionRepo.WithTx(ctx, func(ctx context.Context) (interface{}, error) {
		var terr error

		oldRefreshToken.Revoked = true
		terr = s.sessionRepo.UpdateRefreshTokenCols(ctx, oldRefreshToken, "revoked")
		if terr != nil {
			return nil, terr
		}

		refreshToken, terr = s.createRefreshToken(
			ctx,
			user,
			oldRefreshToken,
			&entities.GrantParams{},
		)
		if terr != nil {
			return nil, terr
		}

		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (s *SessionService) GetCurrentlyActiveRefreshTokenBySessionID(
	ctx context.Context,
	sessionID string,
) (*entities.RefreshToken, error) {
	refreshToken, exists, err := s.sessionRepo.GetCurrentlyActiveRefreshTokenBySessionID(
		ctx,
		sessionID,
	)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.RefreshTokenNotFound)
	}

	return refreshToken, nil
}

func (s *SessionService) grantAuthenticatedUser(
	ctx context.Context,
	user *entities.User,
	params *entities.GrantParams,
) (*entities.RefreshToken, error) {
	return s.createRefreshToken(ctx, user, nil, params)
}

func (s *SessionService) createRefreshToken(
	ctx context.Context,
	user *entities.User,
	oldToken *entities.RefreshToken,
	params *entities.GrantParams,
) (*entities.RefreshToken, error) {
	token, err := crypto.GenerateURLSafeRandomString(entities.TokenLength)
	if err != nil {
		return nil, err
	}

	refreshToken := &entities.RefreshToken{
		ID:     uid.ID(),
		UserID: user.ID,
		Token:  token,
	}
	if oldToken != nil {
		refreshToken.Parent = oldToken.Token
		refreshToken.SessionID = oldToken.SessionID
	}

	if refreshToken.SessionID == "" {
		defaultAAL := entities.Aal1.String()
		session := &entities.Session{
			ID:        uid.ID(),
			UserID:    refreshToken.UserID,
			AAL:       defaultAAL,
			IP:        params.IP,
			UserAgent: params.UserAgent,
		}

		err = s.sessionRepo.AddSession(ctx, session)
		if err != nil {
			return nil, err
		}

		refreshToken.SessionID = session.ID
	}

	err = s.sessionRepo.AddRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (s *SessionService) IssueSignupToken(
	ctx context.Context,
	user *entities.User,
	params *entities.GrantParams,
) (*schemas.CommonTokenResp, error) {
	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(constants.ExpiresInTokenSignup).Unix()

	claims := &entities.CommonTokenClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.ID,
			IssuedAt:  issuedAt.Unix(),
			ExpiresAt: expiresAt,
		},
		Email: user.Email,
		Scope: entities.SignupTokenScope.Ptr(),
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

func (s *SessionService) GenerateAccessToken(
	ctx context.Context,
	user *entities.User,
	sessionID string,
) (string, int64, error) {
	return s.generateAccessToken(ctx, user, sessionID)
}

func (s *SessionService) generateAccessToken(
	ctx context.Context,
	user *entities.User,
	sessionID string,
) (string, int64, error) {
	if sessionID == "" {
		return "", 0, errors.InternalServer(reason.AccessTokenSessionRequired)
	}

	session, exists, err := s.sessionRepo.GetSessionByID(ctx, sessionID, false)
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
			StandardClaims: jwt.StandardClaims{
				Subject:   user.ID,
				IssuedAt:  issuedAt.Unix(),
				ExpiresAt: expiresAt,
			},
			Email: user.Email,
			Scope: entities.AccessTokenScope.Ptr(),
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
