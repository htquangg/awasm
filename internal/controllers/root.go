package controllers

import (
	"github.com/htquangg/a-wasm/internal/services"
)

type Controllers struct {
	Health     *HealthController
	Endpoint   *EndpointController
	Deployment *DeploymentController
	Live       *LiveController
	Preview    *PreviewController
	User       *UserController
	ApiKey     *ApiKeyController
}

func New(services *services.Sevices) *Controllers {
	return &Controllers{
		Health:     NewHealthController(services.Health),
		Endpoint:   NewEndpointController(services.Endpoint),
		Deployment: NewDeploymentController(services.Deployment),
		Live:       NewLiveController(services.Endpoint),
		Preview:    NewPreviewController(services.Deployment),
		User:       NewUserController(services.User),
		ApiKey:     NewApiKeyController(services.ApiKey),
	}
}
