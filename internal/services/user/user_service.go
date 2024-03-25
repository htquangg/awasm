package user

import (
	"context"
	"strings"

	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/pkg/crypto"
	"github.com/htquangg/a-wasm/pkg/uid"

	"github.com/segmentfault/pacman/errors"
)

const (
	TokenLength = 32
)

type (
	UserRepo interface {
		Add(ctx context.Context, user *entities.User) error
		GetUserIDWithEmail(ctx context.Context, email string) (string, bool, error)
	}

	UserAuthRepo interface {
		AddToken(ctx context.Context, userID string, token string, aal string, ip string, userAgent string) error
	}

	UserService struct {
		secretEncryptionKey []byte
		hashingKey          []byte

		userRepo     UserRepo
		userAuthRepo UserAuthRepo
	}
)

func NewUserService(
	userRepo UserRepo,
	userAuthRepo UserAuthRepo,
	secretEncryptionKey []byte,
	hashingKey []byte,
) *UserService {
	return &UserService{
		secretEncryptionKey: secretEncryptionKey,
		hashingKey:          hashingKey,
		userRepo:            userRepo,
		userAuthRepo:        userAuthRepo,
	}
}

func (s *UserService) SignUp(ctx context.Context, req *schemas.SignUpReq) (*schemas.SignUpResp, error) {
	email := strings.ToLower(req.Email)
	// TODO: verify email OTT

	_, exists, err := s.userRepo.GetUserIDWithEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.BadRequest(reason.EmailDuplicate)
	}

	userID, err := s.createUser(ctx, email)
	if err != nil {
		return nil, err
	}

	token, err := crypto.GenerateURLSafeRandomString(TokenLength)
	if err != nil {
		return nil, err
	}
	err = s.userAuthRepo.AddToken(ctx, userID, token, entities.AAL0.String(), req.IP, req.UserAgent)
	if err != nil {
		return nil, err
	}

	return &schemas.SignUpResp{
		Token: token,
	}, err
}

func (s *UserService) createUser(ctx context.Context, email string) (string, error) {
	encryptedEmail, err := crypto.Encrypt(email, s.secretEncryptionKey)
	if err != nil {
		return "", err
	}
	emailHash, err := crypto.GetHash(email, s.hashingKey)
	if err != nil {
		return "", err
	}

	user := &entities.User{}
	user.ID = uid.ID()
	user.EncryptedEmail = encryptedEmail.Cipher
	user.EmailDecryptionNonce = encryptedEmail.Nonce
	user.EmailHash = emailHash

	err = s.userRepo.Add(ctx, user)

	return user.ID, err
}
