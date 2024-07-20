package endpoint_common

import (
	"context"

	"github.com/htquangg/awasm/internal/entities"
)

type EndpointCommonRepo interface {
	GetEndpointByID(ctx context.Context, id string) (*entities.Endpoint, bool, error)
}

type EndpointCommonService struct {
	endpointRepo EndpointCommonRepo
}

func NewEndpointCommonService(endpointRepo EndpointCommonRepo) *EndpointCommonService {
	return &EndpointCommonService{
		endpointRepo: endpointRepo,
	}
}
