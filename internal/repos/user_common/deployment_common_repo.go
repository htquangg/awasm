package user_common

import (
	"context"

	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/a-wasm/internal/base/db"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/services/user_common"
)

type userCommonRepo struct {
	db db.DB
}

func (r *userCommonRepo) GetUserByID(
	ctx context.Context,
	id string,
) (user *entities.User, exists bool, err error) {
	user = &entities.User{}
	exists, err = r.db.Engine(ctx).ID(id).Get(user)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return
}

func NewUserCommonRepo(db db.DB) user_common.UserCommonRepo {
	return &userCommonRepo{
		db: db,
	}
}
