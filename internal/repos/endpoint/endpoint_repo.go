package endpoint

import (
	"context"

	"github.com/htquangg/a-wasm/internal/db"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/reason"
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
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *endpointRepo) GetByID(
	ctx context.Context,
	id string,
) (endpoint *entities.Endpoint, exists bool, err error) {
	endpoint = &entities.Endpoint{}
	exists, err = r.db.Engine(ctx).ID(id).Get(endpoint)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return
}

func (r *endpointRepo) UpdateActiveDeployment(ctx context.Context, endpointID string, deploymentID string) (err error) {
	_, err = r.db.Engine(ctx).ID(endpointID).Update(&entities.Endpoint{ActiveDeploymentID: deploymentID})

	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return
}
