package repos

import (
	"github.com/htquangg/a-wasm/internal/db"
	"github.com/htquangg/a-wasm/internal/repos/deployment"
	"github.com/htquangg/a-wasm/internal/repos/deployment_common"
	"github.com/htquangg/a-wasm/internal/repos/endpoint"
	"github.com/htquangg/a-wasm/internal/repos/endpoint_common"
	deployment_svc "github.com/htquangg/a-wasm/internal/services/deployment"
	deployment_common_svc "github.com/htquangg/a-wasm/internal/services/deployment_common"
	endpoint_svc "github.com/htquangg/a-wasm/internal/services/endpoint"
	envpoint_common_svc "github.com/htquangg/a-wasm/internal/services/endpoint_common"
)

type Repos struct {
	db               db.DB
	Endpoint         endpoint_svc.EndpointRepo
	EndpointCommon   envpoint_common_svc.EndpointCommonRepo
	Deployment       deployment_svc.DeploymentRepo
	DeploymentCommon deployment_common_svc.DeploymentCommonRepo
}

func New(db db.DB) *Repos {
	return &Repos{
		db:               db,
		Endpoint:         endpoint.NewEndpointRepo(db),
		EndpointCommon:   endpoint_common.NewEndpointCommonRepo(db),
		Deployment:       deployment.NewDeploymentRepo(db),
		DeploymentCommon: deployment_common.NewDeploymentCommonRepo(db),
	}
}
