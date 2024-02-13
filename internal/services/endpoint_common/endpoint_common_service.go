package endpoint_common

import (
	"context"

	"github.com/htquangg/a-wasm/internal/entities"
)

type EndpointCommonRepo interface {
	GetByID(ctx context.Context, id string) (*entities.Endpoint, bool, error)
}

type EndpointCommonService struct {
	endpointRepo EndpointCommonRepo
}

func NewEndpointCommonService(endpointRepo EndpointCommonRepo) *EndpointCommonService {
	return &EndpointCommonService{
		endpointRepo: endpointRepo,
	}
}
