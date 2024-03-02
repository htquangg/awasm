package endpoint

import (
	"context"
	"strings"

	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/protocluster"
	"github.com/htquangg/a-wasm/internal/protocluster/grains/messages"
	"github.com/htquangg/a-wasm/internal/reason"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/deployment_common"
	"github.com/htquangg/a-wasm/pkg/uid"

	"github.com/jinzhu/copier"
	"github.com/segmentfault/pacman/errors"
)

type (
	EndpointRepo interface {
		Add(ctx context.Context, endpoint *entities.Endpoint) error
		GetByID(ctx context.Context, id string) (*entities.Endpoint, bool, error)
		UpdateActiveDeployment(ctx context.Context, endpointID, deploymentID string) error
	}

	EndpointService struct {
		protocluster   *protocluster.Cluster
		endpointRepo   EndpointRepo
		deploymentRepo deployment_common.DeploymentCommonRepo
	}
)

func NewEndpointService(endpointRepo EndpointRepo, deploymentRepo deployment_common.DeploymentCommonRepo,
	protoCluster *protocluster.Cluster,
) *EndpointService {
	return &EndpointService{
		protocluster:   protoCluster,
		endpointRepo:   endpointRepo,
		deploymentRepo: deploymentRepo,
	}
}

func (s *EndpointService) Publish(
	ctx context.Context,
	req *schemas.PublishEndpointReq,
) (*schemas.PublishEndpointResp, error) {
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

	currentDeploymentID := endpoint.ActiveDeploymentID
	if strings.EqualFold(currentDeploymentID, deployment.ID) {
		return nil, errors.BadRequest(reason.DeploymentAlreadyActivated)
	}

	err = s.endpointRepo.UpdateActiveDeployment(ctx, endpoint.ID, deployment.ID)
	if err != nil {
		return nil, err
	}

	return &schemas.PublishEndpointResp{
		DeploymentID: deployment.ID,
	}, nil
}

func (s *EndpointService) Add(ctx context.Context, req *schemas.AddEndpointReq) (*schemas.AddEndpointResp, error) {
	endpoint := &entities.Endpoint{}
	_ = copier.Copy(endpoint, req)

	endpoint.ID = uid.ID()

	err := s.endpointRepo.Add(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	resp := &schemas.AddEndpointResp{}
	resp.SetFromEndpoint(endpoint)

	return resp, nil
}

func (s *EndpointService) Serve(
	ctx context.Context,
	req *schemas.ServeEndpointReq,
) (*schemas.ServeEndpointResp, error) {
	endpoint, exists, err := s.endpointRepo.GetByID(ctx, req.EndpointID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.EndpointNotFound)
	}

	if !endpoint.HasActiveDeploy() {
		return nil, errors.BadRequest(reason.EndpointHasNotPublished)
	}

	httpRequest := &messages.HTTPRequest{}
	_ = copier.Copy(httpRequest, req)

	httpRequest.ID = uid.ID()
	httpRequest.DeploymentID = endpoint.ActiveDeploymentID
	httpRequest.Runtime = endpoint.Runtime

	result := s.protocluster.Serve(httpRequest)

	resp := &schemas.ServeEndpointResp{}
	_ = copier.Copy(resp, result)

	return resp, nil
}
