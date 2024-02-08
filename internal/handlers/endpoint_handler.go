package handlers

import "github.com/htquangg/a-wasm/internal/services/endpoint"

type EndpointHandler struct {
	endpointService *endpoint.EndpointService
}

func NewEndpointHandler(endpointService *endpoint.EndpointService) *EndpointHandler {
	return &EndpointHandler{
		endpointService: endpointService,
	}
}
