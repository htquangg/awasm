package services

import (
	"github.com/htquangg/a-wasm/internal/protocluster"
	"github.com/htquangg/a-wasm/internal/repos"
	"github.com/htquangg/a-wasm/internal/services/deployment"
	"github.com/htquangg/a-wasm/internal/services/endpoint"
	"github.com/htquangg/a-wasm/internal/services/health"
	"github.com/htquangg/a-wasm/internal/services/user"
)

type Sevices struct {
	repos        *repos.Repos
	protocluster *protocluster.Cluster
	Health       *health.HealthService
	Endpoint     *endpoint.EndpointService
	Deployment   *deployment.DeploymentService
	User         *user.UserService
}

func New(repos *repos.Repos, cfg *Config, protoCluster *protocluster.Cluster) *Sevices {
	return &Sevices{
		repos:        repos,
		protocluster: protoCluster,
		Health:       health.NewHealthService(),
		Endpoint:     endpoint.NewEndpointService(repos.Endpoint, repos.DeploymentCommon, protoCluster),
		Deployment:   deployment.NewDeploymentService(repos.Deployment, repos.EndpointCommon, protoCluster),
		User:         user.NewUserService(repos.User, repos.UserAuth, cfg.SecretEncryptionKey, cfg.HashingKey),
	}
}
