package services

import (
	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/protocluster"
	"github.com/htquangg/a-wasm/internal/repos"
	"github.com/htquangg/a-wasm/internal/services/auth"
	"github.com/htquangg/a-wasm/internal/services/deployment"
	"github.com/htquangg/a-wasm/internal/services/endpoint"
	"github.com/htquangg/a-wasm/internal/services/health"
	"github.com/htquangg/a-wasm/internal/services/session"
	"github.com/htquangg/a-wasm/internal/services/user"
)

type Sevices struct {
	cfg          *config.Config
	repos        *repos.Repos
	protocluster *protocluster.Cluster
	Health       *health.HealthService
	Endpoint     *endpoint.EndpointService
	Deployment   *deployment.DeploymentService
	Auth         *auth.AuthService
	Session      *session.SessionService
	User         *user.UserService
}

func New(cfg *config.Config, repos *repos.Repos, protoCluster *protocluster.Cluster) *Sevices {
	healthService := health.NewHealthService()
	endpointService := endpoint.NewEndpointService(cfg, repos.Endpoint, repos.DeploymentCommon, protoCluster)
	deploymentService := deployment.NewDeploymentService(cfg, repos.Deployment, repos.EndpointCommon, protoCluster)
	authService := auth.NewAuthService(repos.Auth)
	sessionService := session.NewSessionService(cfg, repos.Session)
	userService := user.NewUserService(cfg, repos.User, sessionService)

	return &Sevices{
		cfg:          cfg,
		repos:        repos,
		protocluster: protoCluster,
		Health:       healthService,
		Endpoint:     endpointService,
		Deployment:   deploymentService,
		Auth:         authService,
		Session:      sessionService,
		User:         userService,
	}
}
