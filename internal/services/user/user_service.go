package user

import (
	"context"
	"strings"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/session"
	"github.com/htquangg/a-wasm/pkg/crypto"
	"github.com/htquangg/a-wasm/pkg/uid"
)

type (
	UserRepo interface {
		AddUser(ctx context.Context, user *entities.User) error
		GetUserByID(ctx context.Context, id string) (*entities.User, bool, error)
		GetUserWithEmail(ctx context.Context, email string) (*entities.User, bool, error)
	}

	UserService struct {
		cfg            *config.Config
		userRepo       UserRepo
		sessionService *session.SessionService
	}
)

func NewUserService(
	cfg *config.Config,
	userRepo UserRepo,
	sessionService *session.SessionService,
) *UserService {
	return &UserService{
		cfg:            cfg,
		userRepo:       userRepo,
		sessionService: sessionService,
	}
}

func (s *UserService) VerifyEmail(ctx context.Context, req *schemas.SignUpReq) (*schemas.AccessTokenResp, error) {
	email := strings.ToLower(req.Email)
	// TODO: verify email OTT

	user, exists, err := s.userRepo.GetUserWithEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	authenticationMethod := entities.OTP
	if !exists {
		authenticationMethod = entities.EmailSignup
		user, err = s.createUser(ctx, email)
		if err != nil {
			return nil, err
		}
	}

	accessTokenResp, err := s.sessionService.IssueRefreshToken(ctx, user, authenticationMethod, &entities.GrantParams{
		IP:        req.IP,
		UserAgent: req.UserAgent,
	})

	return accessTokenResp, err
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
