package user_common

import (
	"context"

	"github.com/htquangg/a-wasm/internal/entities"
)

type UserCommonRepo interface {
	GetUserByID(ctx context.Context, id string) (*entities.User, bool, error)
}

type UserCommonService struct {
	userRepo UserCommonRepo
}

func NewUserCommonService(userRepo UserCommonRepo) *UserCommonService {
	return &UserCommonService{
		userRepo: userRepo,
	}
}
