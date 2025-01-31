package endpoint

import (
	"context"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/awasm/config"
	"github.com/htquangg/awasm/internal/base/reason"
	"github.com/htquangg/awasm/internal/constants"
	"github.com/htquangg/awasm/internal/entities"
	"github.com/htquangg/awasm/internal/protocluster"
	"github.com/htquangg/awasm/internal/protocluster/grains/messages"
	"github.com/htquangg/awasm/internal/schemas"
	"github.com/htquangg/awasm/internal/services/deployment_common"
	"github.com/htquangg/awasm/pkg/uid"
)

type (
	EndpointRepo interface {
		AddEndpoint(ctx context.Context, endpoint *entities.Endpoint) error
		RemoveEndpointByID(ctx context.Context, id string) error
		GetEndpointByID(ctx context.Context, id string) (*entities.Endpoint, bool, error)
		UpdateActiveDeployment(ctx context.Context, endpointID, deploymentID string) error
	}

	EndpointService struct {
		cfg            *config.Config
		protocluster   *protocluster.Cluster
		endpointRepo   EndpointRepo
		deploymentRepo deployment_common.DeploymentCommonRepo
	}
)

func NewEndpointService(
	cfg *config.Config,
	endpointRepo EndpointRepo,
	deploymentRepo deployment_common.DeploymentCommonRepo,
	protoCluster *protocluster.Cluster,
) *EndpointService {
	return &EndpointService{
		cfg:            cfg,
		protocluster:   protoCluster,
		endpointRepo:   endpointRepo,
		deploymentRepo: deploymentRepo,
	}
}

func (s *EndpointService) Publish(
	ctx context.Context,
	req *schemas.PublishEndpointReq,
) (*schemas.PublishEndpointResp, error) {
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
	if endpoint.UserID != req.UserID {
		return nil, errors.Forbidden(reason.EndpointAccessDenied)
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
		IngressURL:   s.cfg.IngressURL + constants.LiveIngressPath + deployment.EndpointID,
	}, nil
}

func (s *EndpointService) AddEndpoint(
	ctx context.Context,
	req *schemas.AddEndpointReq,
) (*schemas.AddEndpointResp, error) {
	endpoint := &entities.Endpoint{}
	_ = copier.Copy(endpoint, req)

	endpoint.ID = uid.ID()

	err := s.endpointRepo.AddEndpoint(ctx, endpoint)
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
	endpoint, exists, err := s.endpointRepo.GetEndpointByID(ctx, req.EndpointID)
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
