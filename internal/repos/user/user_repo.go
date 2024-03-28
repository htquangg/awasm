package user

import (
	"context"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/db"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/services/user"
	"github.com/htquangg/a-wasm/pkg/crypto"

	"github.com/segmentfault/pacman/errors"
)

type userRepo struct {
	cfg *config.Config
	db  db.DB
}

func NewUserRepo(
	cfg *config.Config,
	db db.DB,
) user.UserRepo {
	return &userRepo{
		cfg: cfg,
		db:  db,
	}
}

func (r *userRepo) Add(ctx context.Context, user *entities.User) error {
	_, err := r.db.Engine(ctx).Insert(user)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *userRepo) GetUserIDWithEmail(ctx context.Context, email string) (userID string, exists bool, err error) {
	emailHash, err := crypto.GetHash(email, r.cfg.Key.HashBytes)
	if err != nil {
		return "", false, err
	}

	exists, err = r.db.Engine(ctx).SQL("SELECT id FROM users WHERE email_hash = $1", emailHash).Get(&userID)
	if err != nil {
		return "", false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return userID, exists, err
}
