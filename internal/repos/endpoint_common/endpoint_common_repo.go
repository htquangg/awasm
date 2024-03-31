package endpoint_common

import (
	"context"

	"github.com/htquangg/a-wasm/internal/base/db"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/services/endpoint_common"

	"github.com/segmentfault/pacman/errors"
)

type endpointCommonRepo struct {
	db db.DB
}

func NewEndpointCommonRepo(db db.DB) endpoint_common.EndpointCommonRepo {
	return &endpointCommonRepo{
		db: db,
	}
}

func (r *endpointCommonRepo) GetEndpointByID(
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
