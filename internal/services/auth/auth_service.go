package auth

import (
	"context"

	"github.com/htquangg/a-wasm/internal/entities"
)

type AuthRepo interface {
	GetUserCacheInfo(ctx context.Context, userID string) (userInfo *entities.UserCacheInfo, err error)
	SetUserCacheInfo(ctx context.Context, userID string, userInfo *entities.UserCacheInfo) error
	RemoveUserCacheInfo(ctx context.Context, userID string) (err error)
}

type AuthService struct {
	authRepo AuthRepo
}

func NewAuthService(authRepo AuthRepo) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}
