package handlers

import "github.com/htquangg/a-wasm/internal/services"

type Handlers struct {
	Health *HealthHandler
}

func New(services *services.Sevices) *Handlers {
	return &Handlers{
		Health: NewHealthHandler(services.Health),
	}
}
