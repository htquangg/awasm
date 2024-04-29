package repos

import (
	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/cache"
	"github.com/htquangg/a-wasm/internal/base/db"
	"github.com/htquangg/a-wasm/internal/repos/auth"
	"github.com/htquangg/a-wasm/internal/repos/deployment"
	"github.com/htquangg/a-wasm/internal/repos/deployment_common"
	"github.com/htquangg/a-wasm/internal/repos/endpoint"
	"github.com/htquangg/a-wasm/internal/repos/endpoint_common"
	"github.com/htquangg/a-wasm/internal/repos/health"
	"github.com/htquangg/a-wasm/internal/repos/mailer"
	"github.com/htquangg/a-wasm/internal/repos/session"
	"github.com/htquangg/a-wasm/internal/repos/user"
	auth_svc "github.com/htquangg/a-wasm/internal/services/auth"
	deployment_svc "github.com/htquangg/a-wasm/internal/services/deployment"
	deployment_common_svc "github.com/htquangg/a-wasm/internal/services/deployment_common"
	endpoint_svc "github.com/htquangg/a-wasm/internal/services/endpoint"
	envpoint_common_svc "github.com/htquangg/a-wasm/internal/services/endpoint_common"
	health_svc "github.com/htquangg/a-wasm/internal/services/health"
	mailer_svc "github.com/htquangg/a-wasm/internal/services/mailer"
	session_svc "github.com/htquangg/a-wasm/internal/services/session"
	user_svc "github.com/htquangg/a-wasm/internal/services/user"
)

type Repos struct {
	cfg *config.Config
	db  db.DB

	Health           health_svc.HealthRepo
	Endpoint         endpoint_svc.EndpointRepo
	EndpointCommon   envpoint_common_svc.EndpointCommonRepo
	Deployment       deployment_svc.DeploymentRepo
	DeploymentCommon deployment_common_svc.DeploymentCommonRepo
	Auth             auth_svc.AuthRepo
	Session          session_svc.SessionRepo
	User             user_svc.UserRepo
	UserAuth         user_svc.UserAuthRepo
	Mailer           mailer_svc.MailerRepo
}

func New(cfg *config.Config, db db.DB, cacher cache.Cacher) *Repos {
	return &Repos{
		cfg:              cfg,
		db:               db,
		Health:           health.NewHealthRepo(db, cacher),
		Endpoint:         endpoint.NewEndpointRepo(db),
		EndpointCommon:   endpoint_common.NewEndpointCommonRepo(db),
		Deployment:       deployment.NewDeploymentRepo(db),
		DeploymentCommon: deployment_common.NewDeploymentCommonRepo(db),
		Auth:             auth.NewAuthRepo(db, cacher),
		Session:          session.NewSessionRepo(db),
		User:             user.NewUserRepo(cfg, db),
		UserAuth:         user.NewUserAuthRepo(cfg, db),
		Mailer:           mailer.NewMailerRepo(db, cacher),
	}
}
