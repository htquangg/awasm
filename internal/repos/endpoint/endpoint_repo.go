package endpoint

import (
	"context"

	"github.com/htquangg/a-wasm/internal/db"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/handlers/resp"
	"github.com/htquangg/a-wasm/internal/services/endpoint"
	"github.com/segmentfault/pacman/errors"
)

type endpointRepo struct {
	db db.DB
}

func NewEndpointRepo(db db.DB) endpoint.EndpointRepo {
	return &endpointRepo{
		db: db,
	}
}

func (r *endpointRepo) Add(ctx context.Context, endpoint *entities.Endpoint) error {
	_, err := r.db.Engine(ctx).Insert(endpoint)
	if err != nil {
		return errors.InternalServer(resp.DatabaseError).WithError(err).WithStack()
	}

	return nil
}
