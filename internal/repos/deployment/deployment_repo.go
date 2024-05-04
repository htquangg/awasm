package deployment

import (
	"context"

	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/a-wasm/internal/base/db"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/services/deployment"
)

type deploymentRepo struct {
	db db.DB
}

func NewDeploymentRepo(db db.DB) deployment.DeploymentRepo {
	return &deploymentRepo{
		db: db,
	}
}

func (r *deploymentRepo) AddDeployment(ctx context.Context, deployment *entities.Deployment) error {
	_, err := r.db.Engine(ctx).Insert(deployment)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *deploymentRepo) RemoveDeploymentByID(ctx context.Context, id string) error {
	_, err := r.db.Engine(ctx).ID(id).Delete(new(entities.Deployment))
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *deploymentRepo) GetDeploymentByID(
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
