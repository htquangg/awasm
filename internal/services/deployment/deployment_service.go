package deployment

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/protocluster"
	"github.com/htquangg/a-wasm/internal/protocluster/grains/messages"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/endpoint_common"
	"github.com/htquangg/a-wasm/pkg/uid"
)

type (
	DeploymentRepo interface {
		AddDeployment(ctx context.Context, deployment *entities.Deployment) error
		RemoveDeploymentByID(ctx context.Context, id string) error
		GetDeploymentByID(ctx context.Context, id string) (*entities.Deployment, bool, error)
	}

	DeploymentService struct {
		cfg            *config.Config
		protocluster   *protocluster.Cluster
		deploymentRepo DeploymentRepo
		endpointRepo   endpoint_common.EndpointCommonRepo
	}
)

func NewDeploymentService(
	cfg *config.Config,
	deploymentRepo DeploymentRepo,
	endpointRepo endpoint_common.EndpointCommonRepo,
	protoCluster *protocluster.Cluster,
) *DeploymentService {
	return &DeploymentService{
		cfg:            cfg,
		protocluster:   protoCluster,
		deploymentRepo: deploymentRepo,
		endpointRepo:   endpointRepo,
	}
}

func (s *DeploymentService) AddDeployment(
	ctx context.Context,
	req *schemas.AddDeploymentReq,
) (*schemas.AddDeploymentResp, error) {
	endpoint, exists, err := s.endpointRepo.GetEndpointByID(ctx, req.EndpointID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.EndpointNotFound)
	}
	if endpoint.UserID != req.UserID {
		return nil, errors.Forbidden(reason.EndpointAccessDenied)
	}

	deployment := entities.NewFromEndpoint(endpoint, req.UserID, req.Data)

	err = s.deploymentRepo.AddDeployment(ctx, deployment)
	if err != nil {
		return nil, err
	}

	resp := &schemas.AddDeploymentResp{}
	resp.SetFromDeployment(deployment)
	resp.IngressURL = s.cfg.IngressURL + constants.PreviewIngressPath + resp.ID

	return resp, nil
}

func (s *DeploymentService) Serve(
	ctx context.Context,
	req *schemas.ServeDeploymentReq,
) (*schemas.ServeDeploymentResp, error) {
	deployment, exists, err := s.deploymentRepo.GetDeploymentByID(ctx, req.DeploymentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.DeploymentNotFound)
	}

	canAccess := req.UserID == deployment.UserID
	if !canAccess {
		return nil, errors.Forbidden(reason.DeploymentAccessDenied)
	}

	endpoint, exists, err := s.endpointRepo.GetEndpointByID(ctx, deployment.EndpointID)
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

	resp := &schemas.ServeDeploymentResp{}
	_ = copier.Copy(resp, result)

	return resp, nil
}
