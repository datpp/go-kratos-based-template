package biz

import (
	"context"

	v1 "github.com/datpp/go-kratos-based-template/api/healthcheck/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	Unspecified = errors.New(500, v1.HealthcheckV1ErrorReason_HEALTHCHECK_UNSPECIFIED.String(), "unknown error")
)

// Healthcheck is a Healthcheck model.
type Healthcheck struct {
	Status string `json:"status"`
}

// HealthcheckRepo is a Healthcheck repo.
type HealthcheckRepo interface {
	Save(context.Context, *Healthcheck) (*Healthcheck, error)
	Update(context.Context, *Healthcheck) (*Healthcheck, error)
	FindByID(context.Context, int64) (*Healthcheck, error)
	ListByHello(context.Context, string) ([]*Healthcheck, error)
	ListAll(context.Context) ([]*Healthcheck, error)
}

// HealthcheckUsecase is a Healthcheck usecase.
type HealthcheckUsecase struct {
	repo HealthcheckRepo
	log  *log.Helper
}

// NewHealthcheckUsecase new a Healthcheck usecase.
func NewHealthcheckUsecase(repo HealthcheckRepo, logger log.Logger) *HealthcheckUsecase {
	return &HealthcheckUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateHealthcheck creates a Healthcheck, and returns the new Healthcheck.
func (uc *HealthcheckUsecase) CreateHealthcheck(ctx context.Context, h *Healthcheck) (*Healthcheck, error) {
	uc.log.WithContext(ctx).Infof("CreateHealthcheck: %v", h.Status)
	return uc.repo.Save(ctx, h)
}
