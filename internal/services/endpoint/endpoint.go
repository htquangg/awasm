package endpoint

import "github.com/htquangg/a-wasm/internal/db"

type EndpoinService struct {
	db   db.DB
	repo *endpointRepo
}

func NewEndpointService(db db.DB) *EndpoinService {
	return &EndpoinService{
		db:   db,
		repo: newEndpointRepo(db),
	}
}
