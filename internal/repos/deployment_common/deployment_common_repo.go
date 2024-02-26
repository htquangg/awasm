package deployment_common

import (
	"context"

	"github.com/htquangg/a-wasm/internal/db"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/reason"
	"github.com/htquangg/a-wasm/internal/services/deployment_common"

	"github.com/segmentfault/pacman/errors"
)

type deploymentCommonRepo struct {
	db db.DB
}

func NewDeploymentCommonRepo(db db.DB) deployment_common.DeploymentCommonRepo {
	return &deploymentCommonRepo{
		db: db,
	}
}

func (r *deploymentCommonRepo) GetByID(
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
