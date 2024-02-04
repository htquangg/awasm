package services

import (
	"github.com/htquangg/a-wasm/internal/db"
	"github.com/htquangg/a-wasm/internal/services/endpoint"
)

type Service struct {
	Endpoint *endpoint.Endpoint
}

func New(db db.DB) *Service {
	return &Service{
		Endpoint: endpoint.New(db),
	}
}
