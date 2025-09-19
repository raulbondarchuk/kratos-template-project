//go:build wireinject
// +build wireinject

package main

import (
	"service/internal/broker"
	"service/internal/conf/v1"
	"service/internal/data"
	examplev1 "service/internal/feature/example/v1"
	"service/internal/server"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(app *conf.App, serverConf *conf.Server, dataConf *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		// infra
		server.ProviderSet,
		data.ProviderSet,
		broker.ProviderSet,

		// modules

		examplev1.ProviderSet,

		// single build + distribution to servers
		BuildAllRegistrars,
		ProvideHTTPRegistrers,
		ProvideGRPCRegistrers,

		newApp,
	))
}
