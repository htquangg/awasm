package repos

import (
	"github.com/htquangg/a-wasm/internal/db"
	"github.com/htquangg/a-wasm/internal/repos/deployment"
	"github.com/htquangg/a-wasm/internal/repos/endpoint"
	"github.com/htquangg/a-wasm/internal/repos/endpoint_common"
	deployment_svc "github.com/htquangg/a-wasm/internal/services/deployment"
	endpoint_svc "github.com/htquangg/a-wasm/internal/services/endpoint"
	envpoind_common_svc "github.com/htquangg/a-wasm/internal/services/endpoint_common"
)

type Repos struct {
	db             db.DB
	Endpoint       endpoint_svc.EndpointRepo
	EndpointCommon envpoind_common_svc.EndpointCommonRepo
	Deployment     deployment_svc.DeploymentRepo
}

func New(db db.DB) *Repos {
	return &Repos{
		db:             db,
		Endpoint:       endpoint.NewEndpointRepo(db),
		EndpointCommon: endpoint_common.NewEndpointCommonRepo(db),
		Deployment:     deployment.NewDeploymentRepo(db),
	}
}
