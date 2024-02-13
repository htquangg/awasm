package services

import (
	"github.com/htquangg/a-wasm/internal/repos"
	"github.com/htquangg/a-wasm/internal/services/deployment"
	"github.com/htquangg/a-wasm/internal/services/endpoint"
	"github.com/htquangg/a-wasm/internal/services/health"
)

type Sevices struct {
	repos      *repos.Repos
	Health     *health.HealthService
	Endpoint   *endpoint.EndpointService
	Deployment *deployment.DeploymentService
}

func New(repos *repos.Repos) *Sevices {
	return &Sevices{
		repos:      repos,
		Health:     health.NewHealthService(),
		Endpoint:   endpoint.NewEndpointService(repos.Endpoint),
		Deployment: deployment.NewDeploymentService(repos.Deployment, repos.EndpointCommon),
	}
}
