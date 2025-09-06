//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/datpp/go-kratos-based-template/internal/biz"
	"github.com/datpp/go-kratos-based-template/internal/conf"
	"github.com/datpp/go-kratos-based-template/internal/data"
	"github.com/datpp/go-kratos-based-template/internal/server"
	"github.com/datpp/go-kratos-based-template/internal/service"
	"github.com/datpp/go-kratos-based-template/pkg/utils/bootstrap"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// initApp init kratos application.
func initApp(bootstrap.ServiceInfo, *conf.Server, *conf.Data, *conf.Registry, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
