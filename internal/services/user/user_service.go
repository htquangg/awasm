package user

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/handler"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/mailer"
	"github.com/htquangg/a-wasm/internal/services/session"
	"github.com/htquangg/a-wasm/pkg/converter"
	"github.com/htquangg/a-wasm/pkg/crypto"
	"github.com/htquangg/a-wasm/pkg/uid"
)

const retryLoopDuration = 5 * time.Second

type (
	UserRepo interface {
		WithTx(
			parentCtx context.Context,
			f func(ctx context.Context) (interface{}, error),
		) (interface{}, error)
		AddUser(ctx context.Context, user *entities.User) error
		GetUserByID(ctx context.Context, id string) (*entities.User, bool, error)
		GetUserWithEmail(ctx context.Context, email string) (*entities.User, bool, error)
	}

	UserAuthRepo interface {
		AddSRPChallenge(ctx context.Context, srpChallege *entities.SrpChallenge) error
		AddSrpAuthTemp(ctx context.Context, srpAuthTemp *entities.SrpAuthTemp) error
		CompleteEmailAccountSignup(
			ctx context.Context,
			accountInfo *schemas.CompleteEmailSignupInfo,
		) error
		GetTempSRPSetupByID(ctx context.Context, id string) (*entities.SrpAuthTemp, bool, error)
		GetSRPChallengeByID(ctx context.Context, id string) (*entities.SrpChallenge, bool, error)
		IncrementSrpChallengeAttemptCount(ctx context.Context, id string) error
		SetSrpChallengeVerified(ctx context.Context, challengeID string) error
		GetKeyAttributeWithUserID(
			ctx context.Context,
			userID string,
		) (*entities.KeyAttribute, bool, error)
		GetSRPAttribute(
			ctx context.Context,
			userID string,
		) (*schemas.GetSRPAttributeResp, bool, error)
		GetSRPAuthWithSRPUserID(
			ctx context.Context,
			srpUserID string,
		) (*entities.SrpAuth, bool, error)
	}

	UserService struct {
		cfg            *config.Config
		userRepo       UserRepo
		userAuthRepo   UserAuthRepo
		sessionService *session.SessionService
		mailerService  *mailer.MailerService
	}
)

func NewUserService(
	cfg *config.Config,
	userRepo UserRepo,
	userAuthRepo UserAuthRepo,
	sessionService *session.SessionService,
	mailerSercice *mailer.MailerService,
) *UserService {
	return &UserService{
		cfg:            cfg,
		userRepo:       userRepo,
		userAuthRepo:   userAuthRepo,
		sessionService: sessionService,
		mailerService:  mailerSercice,
	}
}

func (s *UserService) RefreshTokenGrant(
	ctx context.Context,
	req *schemas.RefreshTokenReq,
) (*schemas.RefreshTokenResp, error) {
	retryStart := time.Now()
	retry := true

	for retry && time.Since(retryStart).Seconds() < retryLoopDuration.Seconds() {
		retry = false

		user, refreshToken, session, err := s.sessionService.GetUserWithRefreshToken(
			ctx,
			req.RefreshToken,
			false,
		)
		if err != nil {
			return nil, err
		}

		if session != nil {
			result := session.CheckValidity(
				retryStart,
				&refreshToken.UpdatedAt,
				s.cfg.Session.Timebox,
				s.cfg.Session.InactivityTimeout,
			)

			switch result {
			case entities.SessionValid:
				// do nothing

			case entities.SessionTimedOut:
				return nil, errors.BadRequest(reason.RefreshTokenExpired)

			default:
				return nil, errors.BadRequest(reason.SessionExpired)
			}
		}

		var accessToken string
		var expiresAt int64
		var newAccessTokenResp *schemas.AccessTokenResp

		_, err = s.userRepo.WithTx(ctx, func(ctx context.Context) (interface{}, error) {
			user, refreshToken, session, terr := s.sessionService.GetUserWithRefreshToken(
				ctx,
				req.RefreshToken,
				true,
			)
			if terr != nil {
				if handler.IsNotfoundError(terr) {
					// because forUpdate was set, and the
					// previous check outside the
					// transaction found a user and
					// session, but now we're getting a
					// IsNotFoundError, this means that the
					// user is locked and we need to retry
					// in a few milliseconds
					retry = true
					return nil, terr
				}

				return nil, terr
			}

			var issuedToken *entities.RefreshToken

			if refreshToken.Revoked {
				activeRefreshToken, err := s.sessionService.GetCurrentlyActiveRefreshTokenBySessionID(
					ctx,
					session.ID,
				)
				if err != nil {
					if !handler.IsNotfoundError(err) {
						return nil, err
					}
				}

				if activeRefreshToken != nil && activeRefreshToken.Parent == refreshToken.Token {
					// Token was revoked, but it's the
					// parent of the currently active one.
					// This indicates that the client was
					// not able to store the result when it
					// refreshed token. This case is
					// allowed, provided we return back the
					// active refresh token instead of
					// creating a new one.
					issuedToken = activeRefreshToken
				} else {
					// For a revoked refresh token to be reused, it
					// has to fall within the reuse interval.
					reuseUntil := refreshToken.UpdatedAt.Add(
						time.Second * time.Duration(s.cfg.Security.RefreshTokenReuseInterval))

					if time.Now().After(reuseUntil) {
						return nil, errors.BadRequest(reason.RefreshTokenExpired)
					}
				}
			}

			if issuedToken == nil {
				newToken, terr := s.sessionService.GrantRefreshTokenSwap(ctx, user, refreshToken)
				if terr != nil {
					return nil, terr
				}

				issuedToken = newToken
			}

			accessToken, expiresAt, terr = s.sessionService.GenerateAccessToken(
				ctx,
				user,
				issuedToken.SessionID,
			)
			if terr != nil {
				return nil, terr
			}

			newAccessTokenResp = &schemas.AccessTokenResp{
				CommonTokenResp: schemas.CommonTokenResp{
					AccessToken: accessToken,
					TokenType:   "bearer",
					ExpiresIn:   s.cfg.JWT.Exp,
					ExpiresAt:   expiresAt,
				},
				RefreshToken: issuedToken.Token,
			}

			return nil, nil
		})
		if err != nil {
			if retry && handler.IsNotfoundError(err) {
				time.Sleep(time.Duration(10+rand.Intn(20)) * time.Millisecond)
				continue
			}

			return nil, err
		}

		keyAttribute, exists, err := s.userAuthRepo.GetKeyAttributeWithUserID(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.BadRequest(reason.KeyAttributeNotFound)
		}

		newAccessTokenResp.EncryptedToken, err = crypto.GetEncryptedToken(
			converter.ToB64([]byte(newAccessTokenResp.CommonTokenResp.AccessToken)),
			keyAttribute.PublicKey,
		)
		if err != nil {
			return nil, err
		}
		newAccessTokenResp.CommonTokenResp.AccessToken = ""

		return &schemas.RefreshTokenResp{
			AccessTokenResp: newAccessTokenResp,
		}, nil
	}

	return nil, errors.Conflict(reason.TooManyWrongAttemptsError)
}

func (s *UserService) GetSRPAttribute(
	ctx context.Context,
	req *schemas.GetSRPAttributeReq,
) (*schemas.GetSRPAttributeResp, error) {
	email := strings.ToLower(req.Email)
	user, exists, err := s.userRepo.GetUserWithEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.UserNotFound)
	}

	resp, exists, err := s.userAuthRepo.GetSRPAttribute(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.UserNotFound)
	}

	return resp, nil
}

func (s *UserService) BeginEmailSignupProcess(
	ctx context.Context,
	req *schemas.BeginEmailSignupProcessReq,
) error {
	email := strings.ToLower(req.Email)
	user, exists, err := s.userRepo.GetUserWithEmail(ctx, email)
	if err != nil {
		return err
	}
	if !exists {
		_, err = s.createUser(ctx, email)
		if err != nil {
			return err
		}
	}
	if user.EmailAcceptedAt != nil {
		return errors.BadRequest(reason.EmailDuplicate)
	}

	otp, err := crypto.GenerateOtp()
	if err != nil {
		// OTP generation must always succeed
		panic(err)
	}
	title, body, err := s.mailerService.EmailVerificationTemplate(ctx, otp)
	if err != nil {
		return err
	}

	data := &schemas.EmailCodeContent{
		SourceType: schemas.EmailVerificationSourceType,
		Code:       otp,
		ExpiresAt:  time.Now().UTC().Add(constants.ExpiresInOTPEmailSignup).Unix(),
	}
	err = s.mailerService.SendAndSaveCode(ctx, email, title, body, data)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) VerifyEmailSignup(
	ctx context.Context,
	req *schemas.VerifyEmailSignupReq,
) (*schemas.CommonTokenResp, error) {
	email := strings.ToLower(req.Email)

	check, err := s.mailerService.VerifyCode(
		ctx,
		email,
		schemas.EmailVerificationSourceType,
		req.OTP,
	)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, errors.BadRequest(reason.OTPIncorrect)
	}

	user, exists, err := s.userRepo.GetUserWithEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.EmailNotFound)
	}
	if user.EmailAcceptedAt != nil {
		return nil, errors.BadRequest(reason.EmailDuplicate)
	}

	accessTokenResp, err := s.sessionService.IssueSignupToken(ctx, user, &entities.GrantParams{
		IP:        req.IP,
		UserAgent: req.UserAgent,
	})

	return accessTokenResp, err
}

func (s *UserService) SetupSRPAccountSignup(
	ctx context.Context,
	userID string,
	req *schemas.SetupSRPAccountSignupReq,
) (*schemas.SetupSRPAccountSignupResp, error) {
	srpB, challengeID, err := s.createAndInsertSRPChallenge(
		ctx,
		req.SRPUserID,
		req.SRPVerifier,
		req.SRPA,
	)
	if err != nil {
		return nil, err
	}

	srpAuthTemp := &entities.SrpAuthTemp{}
	srpAuthTemp.ID = uid.ID()
	srpAuthTemp.UserID = userID
	srpAuthTemp.SrpUserID = req.SRPUserID
	srpAuthTemp.Salt = req.SRPSalt
	srpAuthTemp.Verifier = req.SRPVerifier
	srpAuthTemp.SrpChallengeID = challengeID
	err = s.userAuthRepo.AddSrpAuthTemp(ctx, srpAuthTemp)
	if err != nil {
		return nil, err
	}

	return &schemas.SetupSRPAccountSignupResp{
		SetupID: srpAuthTemp.ID,
		SRPB:    *srpB,
	}, nil
}

func (s *UserService) CompleteEmailAccountSignup(
	ctx context.Context,
	user *entities.User,
	req *schemas.CompleteEmailSignupReq,
) (*schemas.CompleteEmailSignupResp, error) {
	tempSRP, exists, err := s.userAuthRepo.GetTempSRPSetupByID(ctx, req.SetupID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.SRPNotFound)
	}

	srpM2, err := s.verifySRPChallenge(ctx, tempSRP.Verifier, tempSRP.SrpChallengeID, req.SRPM1)
	if err != nil {
		return nil, err
	}

	err = s.userAuthRepo.CompleteEmailAccountSignup(ctx, &schemas.CompleteEmailSignupInfo{
		UserID:       user.ID,
		SRPUserID:    tempSRP.SrpUserID,
		Salt:         tempSRP.Salt,
		Verifier:     tempSRP.Verifier,
		KeyAttribute: req.KeyAttribute,
	})
	if err != nil {
		return nil, err
	}

	keyAttribute, exists, err := s.userAuthRepo.GetKeyAttributeWithUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.KeyAttributeNotFound)
	}

	keyAttributeInfo := &schemas.KeyAttributeInfo{}
	keyAttributeInfo.ConvertFromKeyAttributeEntity(keyAttribute)

	return &schemas.CompleteEmailSignupResp{
		KeyAttribute: keyAttributeInfo,
		SRPM2:        *srpM2,
	}, nil
}

func (s *UserService) ChallengeEmailLogin(
	ctx context.Context,
	req *schemas.ChallengeEmailLoginReq,
) (*schemas.ChallengeEmailLoginResp, error) {
	srpAuth, exists, err := s.userAuthRepo.GetSRPAuthWithSRPUserID(ctx, req.SRPUserID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.SRPNotFound)
	}

	srpB, challengeID, err := s.createAndInsertSRPChallenge(
		ctx,
		req.SRPUserID,
		srpAuth.Verifier,
		req.SRPA,
	)
	if err != nil {
		return nil, err
	}

	return &schemas.ChallengeEmailLoginResp{
		ChallengeID: challengeID,
		SRPB:        *srpB,
	}, nil
}

func (s *UserService) VerifyEmailLogin(
	ctx context.Context,
	req *schemas.VerifyEmailLoginReq,
) (*schemas.VerifyEmailLoginResp, error) {
	srpAuth, exists, err := s.userAuthRepo.GetSRPAuthWithSRPUserID(ctx, req.SRPUserID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.SRPNotFound)
	}

	srpM2, err := s.verifySRPChallenge(ctx, srpAuth.Verifier, req.ChallengeID, req.SRPM1)
	if err != nil {
		return nil, err
	}

	user, exists, err := s.userRepo.GetUserByID(ctx, srpAuth.UserID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.UserNotFound)
	}

	keyAttribute, exists, err := s.userAuthRepo.GetKeyAttributeWithUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.KeyAttributeNotFound)
	}

	accessTokenResp, err := s.sessionService.IssueRefreshToken(
		ctx,
		user,
		entities.EmailSignup,
		&entities.GrantParams{},
	)
	if err != nil {
		return nil, err
	}
	accessTokenResp.EncryptedToken, err = crypto.GetEncryptedToken(
		converter.ToB64([]byte(accessTokenResp.CommonTokenResp.AccessToken)),
		keyAttribute.PublicKey,
	)
	if err != nil {
		return nil, err
	}
	accessTokenResp.CommonTokenResp.AccessToken = ""

	keyAttributeInfo := &schemas.KeyAttributeInfo{}
	keyAttributeInfo.ConvertFromKeyAttributeEntity(keyAttribute)

	return &schemas.VerifyEmailLoginResp{
		AccessTokenResp: accessTokenResp,
		KeyAttribute:    keyAttributeInfo,
		SRPM2:           *srpM2,
	}, nil
}

func (s *UserService) createUser(ctx context.Context, email string) (*entities.User, error) {
	encryptedEmail, err := crypto.Encrypt(email, s.cfg.Key.EncryptionBytes)
	if err != nil {
		return nil, err
	}
	emailHash, err := crypto.GetHash(email, s.cfg.Key.HashBytes)
	if err != nil {
		return nil, err
	}

	user := &entities.User{}
	user.ID = uid.ID()
	user.EncryptedEmail = encryptedEmail.Cipher
	user.EmailDecryptionNonce = encryptedEmail.Nonce
	user.EmailHash = emailHash
	err = s.userRepo.AddUser(ctx, user)

	return user, err
}
