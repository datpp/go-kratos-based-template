package data

import (
	"context"

	"github.com/datpp/go-kratos-based-template/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type healthcheckRepo struct {
	data *Data
	log  *log.Helper
}

// NewHealthcheckRepo .
func NewHealthcheckRepo(data *Data, logger log.Logger) biz.HealthcheckRepo {
	return &healthcheckRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *healthcheckRepo) Save(ctx context.Context, h *biz.Healthcheck) (*biz.Healthcheck, error) {
	return h, nil
}

func (r *healthcheckRepo) Update(ctx context.Context, h *biz.Healthcheck) (*biz.Healthcheck, error) {
	return h, nil
}

func (r *healthcheckRepo) FindByID(context.Context, int64) (*biz.Healthcheck, error) {
	return nil, nil
}

func (r *healthcheckRepo) ListByHello(context.Context, string) ([]*biz.Healthcheck, error) {
	return nil, nil
}

func (r *healthcheckRepo) ListAll(context.Context) ([]*biz.Healthcheck, error) {
	return nil, nil
}
