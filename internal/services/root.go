package services

import (
	"github.com/htquangg/a-wasm/internal/protocluster"
	"github.com/htquangg/a-wasm/internal/repos"
	"github.com/htquangg/a-wasm/internal/services/deployment"
	"github.com/htquangg/a-wasm/internal/services/endpoint"
	"github.com/htquangg/a-wasm/internal/services/health"
)

type Sevices struct {
	repos        *repos.Repos
	protocluster *protocluster.Cluster
	Health       *health.HealthService
	Endpoint     *endpoint.EndpointService
	Deployment   *deployment.DeploymentService
}

func New(repos *repos.Repos, protoCluster *protocluster.Cluster) *Sevices {
	return &Sevices{
		repos:        repos,
		protocluster: protoCluster,
		Health:       health.NewHealthService(),
		Endpoint:     endpoint.NewEndpointService(repos.Endpoint, repos.DeploymentCommon, protoCluster),
		Deployment:   deployment.NewDeploymentService(repos.Deployment, repos.EndpointCommon, protoCluster),
	}
}
