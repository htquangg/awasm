package user

import (
	"context"

	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/awasm/config"
	"github.com/htquangg/awasm/internal/base/db"
	"github.com/htquangg/awasm/internal/base/reason"
	"github.com/htquangg/awasm/internal/entities"
	"github.com/htquangg/awasm/internal/services/user"
	"github.com/htquangg/awasm/pkg/crypto"
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

func (r *userRepo) WithTx(
	parentCtx context.Context,
	f func(ctx context.Context) (interface{}, error),
) (interface{}, error) {
	return r.db.WithTx(parentCtx, f)
}

func (r *userRepo) AddUser(ctx context.Context, user *entities.User) (err error) {
	_, err = r.db.Engine(ctx).Insert(user)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *userRepo) GetUserWithEmail(
	ctx context.Context,
	email string,
) (user *entities.User, exists bool, err error) {
	var emailHash string
	emailHash, err = crypto.GetHash(email, r.cfg.Key.HashBytes)
	if err != nil {
		return nil, false, err
	}

	user = &entities.User{}
	exists, err = r.db.Engine(ctx).Where("email_hash = $1", emailHash).Get(user)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	user.Email = email

	return user, exists, err
}

func (r *userRepo) GetUserByID(
	ctx context.Context,
	id string,
) (user *entities.User, exists bool, err error) {
	user = &entities.User{}
	exists, err = r.db.Engine(ctx).Where("id = $1", id).Get(user)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	email, err := crypto.Decrypt(
		user.EncryptedEmail,
		r.cfg.Key.EncryptionBytes,
		user.EmailDecryptionNonce,
	)
	if err != nil {
		return nil, false, err
	}
	user.Email = email

	return user, exists, err
}
