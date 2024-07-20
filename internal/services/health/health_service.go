package health

import (
	"context"

	"github.com/htquangg/awasm/internal/schemas"
)

type HealthRepo interface {
	Check(ctx context.Context) error
}

type HealthService struct {
	healthRepo HealthRepo
}

func NewHealthService(healthRepo HealthRepo) *HealthService {
	return &HealthService{
		healthRepo: healthRepo,
	}
}

func (s *HealthService) CheckHealth(ctx context.Context) (*schemas.CheckHealthResp, error) {
	if err := s.healthRepo.Check(ctx); err != nil {
		return nil, err
	}

	return &schemas.CheckHealthResp{Msg: "API is live!!!"}, nil
}
