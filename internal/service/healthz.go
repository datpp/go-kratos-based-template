package service

import (
	"context"

	v1 "github.com/datpp/go-kratos-based-template/api/healthcheck/v1"
	"github.com/datpp/go-kratos-based-template/internal/biz"
)

// HealthcheckService is a healthcheck service.
type HealthcheckService struct {
	v1.UnimplementedHealthCheckServer

	uc *biz.HealthcheckUsecase
}

// NewHealthcheckService new a healthcheck service.
func NewHealthcheckService(uc *biz.HealthcheckUsecase) *HealthcheckService {
	return &HealthcheckService{uc: uc}
}

// Healthcheck implements healthcheck.HealthcheckServer.
func (s *HealthcheckService) HealthCheck(ctx context.Context, in *v1.HealthCheckRequest) (*v1.HealthCheckResponse, error) {
	h, err := s.uc.CreateHealthcheck(ctx, &biz.Healthcheck{Status: "ok"})
	if err != nil {
		return nil, err
	}
	return &v1.HealthCheckResponse{Status: h.Status}, nil
}
