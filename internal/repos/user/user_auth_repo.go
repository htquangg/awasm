package user

import (
	"context"
	"time"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/db"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/services/user"

	"github.com/segmentfault/pacman/errors"
)

type userAuthRepo struct {
	cfg *config.Config
	db  db.DB
}

func NewUserAuthRepo(
	cfg *config.Config,
	db db.DB,
) user.UserAuthRepo {
	return &userAuthRepo{
		cfg: cfg,
		db:  db,
	}
}

func (r *userAuthRepo) AddToken(
	ctx context.Context,
	userID string,
	token string,
	aal string,
	ip string,
	userAgent string,
) error {
	now := time.Now()

	_, err := r.db.Engine(ctx).Insert(&entities.Token{
		LastUsedAt: &now,
		UserID:     userID,
		Token:      token,
		AAL:        aal,
		IP:         ip,
		UserAgent:  userAgent,
	})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}
