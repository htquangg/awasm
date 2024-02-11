package endpoint

import (
	"context"

	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/schemas"

	"github.com/jinzhu/copier"
)

type (
	EndpointRepo interface {
		Add(ctx context.Context, endpoint *entities.Endpoint) error
	}

	EndpointService struct {
		endpointRepo EndpointRepo
	}
)

func NewEndpointService(endpointRepo EndpointRepo) *EndpointService {
	return &EndpointService{
		endpointRepo: endpointRepo,
	}
}

func (r *EndpointService) Add(ctx context.Context, req *schemas.AddEndpointReq) (*schemas.AddEndpointResp, error) {
	endpoint := &entities.Endpoint{}
	_ = copier.Copy(endpoint, req)

	err := r.endpointRepo.Add(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	result := &schemas.AddEndpointResp{}
	result.SetFromEndpoint(endpoint)

	return result, nil
}
