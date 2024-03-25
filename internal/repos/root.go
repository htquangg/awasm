package repos

import (
	"github.com/htquangg/a-wasm/internal/base/db"
	"github.com/htquangg/a-wasm/internal/repos/deployment"
	"github.com/htquangg/a-wasm/internal/repos/deployment_common"
	"github.com/htquangg/a-wasm/internal/repos/endpoint"
	"github.com/htquangg/a-wasm/internal/repos/endpoint_common"
	"github.com/htquangg/a-wasm/internal/repos/user"
	deployment_svc "github.com/htquangg/a-wasm/internal/services/deployment"
	deployment_common_svc "github.com/htquangg/a-wasm/internal/services/deployment_common"
	endpoint_svc "github.com/htquangg/a-wasm/internal/services/endpoint"
	envpoint_common_svc "github.com/htquangg/a-wasm/internal/services/endpoint_common"
	user_svc "github.com/htquangg/a-wasm/internal/services/user"
)

type Repos struct {
	db               db.DB
	Endpoint         endpoint_svc.EndpointRepo
	EndpointCommon   envpoint_common_svc.EndpointCommonRepo
	Deployment       deployment_svc.DeploymentRepo
	DeploymentCommon deployment_common_svc.DeploymentCommonRepo
	User             user_svc.UserRepo
	UserAuth         user_svc.UserAuthRepo
}

func New(db db.DB, cfg *Config) *Repos {
	return &Repos{
		db:               db,
		Endpoint:         endpoint.NewEndpointRepo(db),
		EndpointCommon:   endpoint_common.NewEndpointCommonRepo(db),
		Deployment:       deployment.NewDeploymentRepo(db),
		DeploymentCommon: deployment_common.NewDeploymentCommonRepo(db),
		User:             user.NewUserRepo(db, cfg.SecretEncryptionKey, cfg.HashingKey),
		UserAuth:         user.NewUserAuthRepo(db, cfg.SecretEncryptionKey, cfg.HashingKey),
	}
}
