package deployment

import (
	"context"

	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/protocluster"
	"github.com/htquangg/a-wasm/internal/protocluster/grains/messages"
	"github.com/htquangg/a-wasm/internal/reason"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/endpoint_common"
	"github.com/htquangg/a-wasm/pkg/uid"

	"github.com/jinzhu/copier"
	"github.com/segmentfault/pacman/errors"
)

type (
	DeploymentRepo interface {
		Add(ctx context.Context, deployment *entities.Deployment) error
		GetByID(ctx context.Context, id string) (*entities.Deployment, bool, error)
	}

	DeploymentService struct {
		protocluster   *protocluster.Cluster
		deploymentRepo DeploymentRepo
		endpointRepo   endpoint_common.EndpointCommonRepo
	}
)

func NewDeploymentService(
	deploymentRepo DeploymentRepo,
	endpointRepo endpoint_common.EndpointCommonRepo,
	protoCluster *protocluster.Cluster,
) *DeploymentService {
	return &DeploymentService{
		protocluster:   protoCluster,
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

func (s *DeploymentService) Serve(
	ctx context.Context,
	req *schemas.ServePreviewReq,
) (*schemas.ServePreviewResp, error) {
	deployment, exists, err := s.deploymentRepo.GetByID(ctx, req.DeploymentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.DeploymentNotFound)
	}

	endpoint, exists, err := s.endpointRepo.GetByID(ctx, deployment.EndpointID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.EndpointNotFound)
	}

	httpRequest := &messages.HTTPRequest{}
	_ = copier.Copy(httpRequest, req)

	httpRequest.ID = uid.ID()
	httpRequest.EndpointID = endpoint.ID
	httpRequest.Runtime = endpoint.Runtime

	result := s.protocluster.Serve(httpRequest)

	resp := &schemas.ServePreviewResp{}
	_ = copier.Copy(resp, result)

	return resp, nil
}
