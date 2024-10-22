package health

import (
	"context"

	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/awasm/internal/base/cache"
	"github.com/htquangg/awasm/internal/base/db"
	"github.com/htquangg/awasm/internal/base/reason"
	"github.com/htquangg/awasm/internal/services/health"
)

type healthRepo struct {
	db     db.DB
	cacher cache.Cacher
}

func NewHealthRepo(db db.DB, cacher cache.Cacher) health.HealthRepo {
	return &healthRepo{
		db:     db,
		cacher: cacher,
	}
}

func (r *healthRepo) Check(ctx context.Context) error {
	if err := r.db.Engine(ctx).Ping(); err != nil {
		return errors.ServiceUnavailable(reason.DatabaseError).WithError(err)
	}

	if err := r.cacher.Ping(ctx); err != nil {
		return errors.ServiceUnavailable(reason.DatabaseError).WithError(err)
	}

	return nil
}
