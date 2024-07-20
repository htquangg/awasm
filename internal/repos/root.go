package repos

import (
	"github.com/htquangg/awasm/config"
	"github.com/htquangg/awasm/internal/base/cache"
	"github.com/htquangg/awasm/internal/base/db"
	"github.com/htquangg/awasm/internal/repos/api_key"
	"github.com/htquangg/awasm/internal/repos/auth"
	"github.com/htquangg/awasm/internal/repos/deployment"
	"github.com/htquangg/awasm/internal/repos/deployment_common"
	"github.com/htquangg/awasm/internal/repos/endpoint"
	"github.com/htquangg/awasm/internal/repos/endpoint_common"
	"github.com/htquangg/awasm/internal/repos/health"
	"github.com/htquangg/awasm/internal/repos/mailer"
	"github.com/htquangg/awasm/internal/repos/session"
	"github.com/htquangg/awasm/internal/repos/user"
	"github.com/htquangg/awasm/internal/repos/user_common"
	api_key_svc "github.com/htquangg/awasm/internal/services/api_key"
	auth_svc "github.com/htquangg/awasm/internal/services/auth"
	deployment_svc "github.com/htquangg/awasm/internal/services/deployment"
	deployment_common_svc "github.com/htquangg/awasm/internal/services/deployment_common"
	endpoint_svc "github.com/htquangg/awasm/internal/services/endpoint"
	envpoint_common_svc "github.com/htquangg/awasm/internal/services/endpoint_common"
	health_svc "github.com/htquangg/awasm/internal/services/health"
	mailer_svc "github.com/htquangg/awasm/internal/services/mailer"
	session_svc "github.com/htquangg/awasm/internal/services/session"
	user_svc "github.com/htquangg/awasm/internal/services/user"
	user_common_svc "github.com/htquangg/awasm/internal/services/user_common"
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
	UserCommon       user_common_svc.UserCommonRepo
	UserAuth         user_svc.UserAuthRepo
	Mailer           mailer_svc.MailerRepo
	ApiKey           api_key_svc.ApiKeyRepo
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
		UserCommon:       user_common.NewUserCommonRepo(db),
		UserAuth:         user.NewUserAuthRepo(cfg, db),
		Mailer:           mailer.NewMailerRepo(db, cacher),
		ApiKey:           api_key.NewApiKeyRepo(db),
	}
}
