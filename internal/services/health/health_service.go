package health

import "github.com/htquangg/a-wasm/internal/schemas"

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) CheckHealth() (*schemas.CheckHealthResp, error) {
	return &schemas.CheckHealthResp{Msg: "API is live!!!"}, nil
}
