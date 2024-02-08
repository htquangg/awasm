package endpoint

import (
	"github.com/htquangg/a-wasm/internal/db"
	"github.com/htquangg/a-wasm/internal/services/endpoint"
)

type endpointRepo struct {
	db db.DB
}

func NewEndpointRepo(db db.DB) endpoint.EndpointRepo {
	return &endpointRepo{
		db: db,
	}
}
