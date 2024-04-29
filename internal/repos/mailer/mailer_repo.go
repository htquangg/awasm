package mailer

import (
	"context"
	"time"

	"github.com/htquangg/a-wasm/internal/base/cache"
	"github.com/htquangg/a-wasm/internal/base/db"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/services/mailer"

	"github.com/segmentfault/pacman/errors"
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
	err := r.cacher.Set(ctx, code, []byte(content), duration)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *mailerRepo) GetCode(ctx context.Context, code string) (string, error) {
	content, exist, err := r.cacher.Get(ctx, code)
	if err != nil {
		return "", errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if !exist {
		return "", nil
	}

	return string(content), nil
}

func (r *mailerRepo) DeleteCode(ctx context.Context, code string) error {
	err := r.cacher.Delete(ctx, code)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}
