package repos

import "github.com/htquangg/a-wasm/internal/db"

type Repos struct {
	db         db.DB
	Deployment *DeploymentRepo
}

func New(db db.DB) *Repos {
	return &Repos{
		db:         db,
		Deployment: NewDeploymentRepo(db),
	}
}
