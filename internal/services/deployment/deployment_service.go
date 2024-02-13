package deployment

import (
	"context"

	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/reason"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/endpoint_common"
	"github.com/segmentfault/pacman/errors"
)

type (
	DeploymentRepo interface {
		Add(ctx context.Context, deployment *entities.Deployment) error
		GetByID(ctx context.Context, id string) (*entities.Deployment, bool, error)
	}

	DeploymentService struct {
		deploymentRepo DeploymentRepo
		endpointRepo   endpoint_common.EndpointCommonRepo
	}
)

func NewDeploymentService(
	deploymentRepo DeploymentRepo,
	endpointRepo endpoint_common.EndpointCommonRepo,
) *DeploymentService {
	return &DeploymentService{
		deploymentRepo: deploymentRepo,
		endpointRepo:   endpointRepo,
	}
}

func (s *DeploymentService) Add(
	ctx context.Context,
	req *schemas.AddDeploymentReq,
) (*schemas.AddDeploymentResp, error) {
	endpoint, exists, err := s.endpointRepo.GetByID(ctx, req.EndpointID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.EndpointNotFound)
	}

	deployment := entities.NewFromEndpoint(endpoint, req.Data)

	err = s.deploymentRepo.Add(ctx, deployment)
	if err != nil {
		return nil, err
	}

	resp := &schemas.AddDeploymentResp{}
	resp.SetFromDeployment(deployment)

	return resp, nil
}
