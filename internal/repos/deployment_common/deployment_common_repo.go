package deployment_common

import (
	"context"

	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/awasm/internal/base/db"
	"github.com/htquangg/awasm/internal/base/reason"
	"github.com/htquangg/awasm/internal/entities"
	"github.com/htquangg/awasm/internal/services/deployment_common"
)

type deploymentCommonRepo struct {
	db db.DB
}

func NewDeploymentCommonRepo(db db.DB) deployment_common.DeploymentCommonRepo {
	return &deploymentCommonRepo{
		db: db,
	}
}

func (r *deploymentCommonRepo) GetDeploymentByID(
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
