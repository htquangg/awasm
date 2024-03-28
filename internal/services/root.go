package services

import (
	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/protocluster"
	"github.com/htquangg/a-wasm/internal/repos"
	"github.com/htquangg/a-wasm/internal/services/deployment"
	"github.com/htquangg/a-wasm/internal/services/endpoint"
	"github.com/htquangg/a-wasm/internal/services/health"
	"github.com/htquangg/a-wasm/internal/services/user"
)

type Sevices struct {
	cfg          *config.Config
	repos        *repos.Repos
	protocluster *protocluster.Cluster
	Health       *health.HealthService
	Endpoint     *endpoint.EndpointService
	Deployment   *deployment.DeploymentService
	User         *user.UserService
}

func New(cfg *config.Config, repos *repos.Repos, protoCluster *protocluster.Cluster) *Sevices {
	return &Sevices{
		cfg:          cfg,
		repos:        repos,
		protocluster: protoCluster,
		Health:       health.NewHealthService(),
		Endpoint:     endpoint.NewEndpointService(repos.Endpoint, repos.DeploymentCommon, protoCluster),
		Deployment:   deployment.NewDeploymentService(repos.Deployment, repos.EndpointCommon, protoCluster),
		User:         user.NewUserService(cfg, repos.User, repos.UserAuth),
	}
}
