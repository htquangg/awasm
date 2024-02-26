package endpoint

import (
	"context"

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
		endpointRepo:   endpointRepo,
		deploymentRepo: deploymentRepo,
	}
}

func (r *EndpointService) Add(ctx context.Context, req *schemas.AddEndpointReq) (*schemas.AddEndpointResp, error) {
	endpoint := &entities.Endpoint{}
	_ = copier.Copy(endpoint, req)

	endpoint.ID = uid.ID()

	err := r.endpointRepo.Add(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	resp := &schemas.AddEndpointResp{}
	resp.SetFromEndpoint(endpoint)

	return resp, nil
}

func (s *EndpointService) Serve(
	ctx context.Context,
	req *schemas.ServeLiveReq,
) (*schemas.ServeLiveResp, error) {
	endpoint, exists, err := s.endpointRepo.GetByID(ctx, req.EndpointID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.DeploymentNotFound)
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

	resp := &schemas.ServeLiveResp{}
	_ = copier.Copy(resp, result)

	return resp, nil
}
