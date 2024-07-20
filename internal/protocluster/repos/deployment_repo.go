package repos

import (
	"context"

	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/awasm/internal/base/db"
	"github.com/htquangg/awasm/internal/base/reason"
	"github.com/htquangg/awasm/internal/entities"
)

type DeploymentRepo struct {
	db db.DB
}

func NewDeploymentRepo(db db.DB) *DeploymentRepo {
	return &DeploymentRepo{
		db: db,
	}
}

func (r *DeploymentRepo) GetByID(
	ctx context.Context,
	id string,
) (deployment *entities.Deployment, exists bool, err error) {
	deployment = &entities.Deployment{}
	exists, err = r.db.Engine(ctx).ID(id).Get(deployment)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return
}
