package mailer

import (
	"context"
	"time"

	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/a-wasm/internal/base/cache"
	"github.com/htquangg/a-wasm/internal/base/db"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/internal/services/mailer"
)

type mailerRepo struct {
	db     db.DB
	cacher cache.Cacher
}

func NewMailerRepo(db db.DB, cacher cache.Cacher) mailer.MailerRepo {
	return &mailerRepo{
		db:     db,
		cacher: cacher,
	}
}

func (r *mailerRepo) SetCode(ctx context.Context, code string, content string, duration time.Duration) error {
	cacheKey := &cache.Key{
		Namespace: constants.MailerOTPNameSpaceCache,
		Key:       code,
	}
	err := r.cacher.Set(ctx, cacheKey, []byte(content), duration)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *mailerRepo) GetCode(ctx context.Context, code string) (string, error) {
	cacheKey := &cache.Key{
		Namespace: constants.MailerOTPNameSpaceCache,
		Key:       code,
	}
	content, exist, err := r.cacher.Get(ctx, cacheKey)
	if err != nil {
		return "", errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if !exist {
		return "", nil
	}

	return string(content), nil
}

func (r *mailerRepo) DeleteCode(ctx context.Context, code string) error {
	cacheKey := &cache.Key{
		Namespace: constants.MailerOTPNameSpaceCache,
		Key:       code,
	}
	err := r.cacher.Delete(ctx, cacheKey)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}
