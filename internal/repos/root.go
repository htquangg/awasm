package repos

import (
	"github.com/htquangg/a-wasm/internal/db"
	"github.com/htquangg/a-wasm/internal/repos/endpoint"
	endpoint_svc "github.com/htquangg/a-wasm/internal/services/endpoint"
)

type Repos struct {
	db       db.DB
	Endpoint endpoint_svc.EndpointRepo
}

func New(db db.DB) *Repos {
	return &Repos{
		db:       db,
		Endpoint: endpoint.NewEndpointRepo(db),
	}
}
