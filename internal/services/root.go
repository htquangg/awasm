package services

import (
	"github.com/htquangg/a-wasm/internal/db"
	"github.com/htquangg/a-wasm/internal/services/endpoint"
)

type Service struct {
	Endpoint *endpoint.EndpoinService
}

func New(db db.DB) *Service {
	return &Service{
		Endpoint: endpoint.NewEndpointService(db),
	}
}
