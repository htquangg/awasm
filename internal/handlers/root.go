package handlers

import "github.com/htquangg/a-wasm/internal/services"

type Handlers struct {
	Health   *HealthHandler
	Endpoint *EndpointHandler
}

func New(services *services.Sevices) *Handlers {
	return &Handlers{
		Health:   NewHealthHandler(services.Health),
		Endpoint: NewEndpointHandler(services.Endpoint),
	}
}
