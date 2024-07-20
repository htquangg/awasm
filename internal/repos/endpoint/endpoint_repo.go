package endpoint

import (
	"context"

	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/awasm/internal/base/db"
	"github.com/htquangg/awasm/internal/base/reason"
	"github.com/htquangg/awasm/internal/entities"
	"github.com/htquangg/awasm/internal/services/endpoint"
)

type endpointRepo struct {
	db db.DB
}

func NewEndpointRepo(db db.DB) endpoint.EndpointRepo {
	return &endpointRepo{
		db: db,
	}
}

func (r *endpointRepo) AddEndpoint(ctx context.Context, endpoint *entities.Endpoint) error {
	_, err := r.db.Engine(ctx).Insert(endpoint)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *endpointRepo) RemoveEndpointByID(ctx context.Context, id string) error {
	_, err := r.db.Engine(ctx).ID(id).Delete(new(entities.Endpoint))
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *endpointRepo) GetEndpointByID(
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

func (r *endpointRepo) UpdateActiveDeployment(
	ctx context.Context,
	endpointID string,
	deploymentID string,
) (err error) {
	_, err = r.db.Engine(ctx).
		ID(endpointID).
		Update(&entities.Endpoint{ActiveDeploymentID: deploymentID})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return
}
