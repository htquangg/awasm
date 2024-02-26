package deployment_common

import (
	"context"

	"github.com/htquangg/a-wasm/internal/entities"
)

type DeploymentCommonRepo interface {
	GetByID(ctx context.Context, id string) (*entities.Deployment, bool, error)
}

type DeploymentCommonService struct {
	deploymentRepo DeploymentCommonRepo
}

func NewDeploymentCommonService(deploymentRepo DeploymentCommonRepo) *DeploymentCommonService {
	return &DeploymentCommonService{
		deploymentRepo: deploymentRepo,
	}
}
