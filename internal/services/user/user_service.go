package user

import (
	"context"
	"strings"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/session"
	"github.com/htquangg/a-wasm/pkg/converter"
	"github.com/htquangg/a-wasm/pkg/crypto"
	"github.com/htquangg/a-wasm/pkg/uid"

	"github.com/segmentfault/pacman/errors"
)

type (
	UserRepo interface {
		AddUser(ctx context.Context, user *entities.User) error
		GetUserByID(ctx context.Context, id string) (*entities.User, bool, error)
		GetUserWithEmail(ctx context.Context, email string) (*entities.User, bool, error)
	}

	UserAuthRepo interface {
		AddSRPChallenge(ctx context.Context, srpChallege *entities.SrpChallenge) error
		AddSrpAuthTemp(ctx context.Context, srpAuthTemp *entities.SrpAuthTemp) error
		CompleteEmailAccountSignup(ctx context.Context, accountInfo *schemas.CompleteEmailSignupInfo) error
		GetTempSRPSetupByID(ctx context.Context, id string) (*entities.SrpAuthTemp, bool, error)
		GetSRPChallengeByID(ctx context.Context, id string) (*entities.SrpChallenge, bool, error)
		IncrementSrpChallengeAttemptCount(ctx context.Context, id string) error
		SetSrpChallengeVerified(ctx context.Context, challengeID string) error
		GetKeyAttributeWithUserID(ctx context.Context, userID string) (*entities.KeyAttribute, bool, error)
		GetSRPAttribute(ctx context.Context, userID string) (*schemas.GetSRPAttributeResp, bool, error)
		GetSRPAuthWithSRPUserID(ctx context.Context, srpUserID string) (*entities.SrpAuth, bool, error)
	}

	UserService struct {
		cfg            *config.Config
		userRepo       UserRepo
		userAuthRepo   UserAuthRepo
		sessionService *session.SessionService
	}
)

func NewUserService(
	cfg *config.Config,
	userRepo UserRepo,
	userAuthRepo UserAuthRepo,
	sessionService *session.SessionService,
) *UserService {
	return &UserService{
		cfg:            cfg,
		userRepo:       userRepo,
		userAuthRepo:   userAuthRepo,
		sessionService: sessionService,
	}
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

func (s *UserService) BeginEmailSignupProcess(ctx context.Context, req *schemas.BeginEmailSignupProcessReq) error {
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

	// TODO: send email with OTP
	return nil
}

func (s *UserService) VerifyEmailSignup(
	ctx context.Context,
	req *schemas.VerifyEmailSignupReq,
) (*schemas.CommonTokenResp, error) {
	email := strings.ToLower(req.Email)
	// TODO: verify email OTP

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
	srpB, challengeID, err := s.createAndInsertSRPChallenge(ctx, req.SRPUserID, req.SRPVerifier, req.SRPA)
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

	accessTokenResp, err := s.sessionService.IssueRefreshToken(ctx, user, entities.EmailSignup, &entities.GrantParams{
		IP:        req.IP,
		UserAgent: req.UserAgent,
	})
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

	return &schemas.CompleteEmailSignupResp{
		AccessTokenResp: accessTokenResp,
		KeyAttribute:    keyAttributeInfo,
		SRPM2:           *srpM2,
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

	srpB, challengeID, err := s.createAndInsertSRPChallenge(ctx, req.SRPUserID, srpAuth.Verifier, req.SRPA)
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

	accessTokenResp, err := s.sessionService.IssueRefreshToken(ctx, user, entities.EmailSignup, &entities.GrantParams{})
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
