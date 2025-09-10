//go:build wireinject
// +build wireinject

package main

import (
	"service/internal/broker"
	"service/internal/conf"
	"service/internal/data"
	"service/internal/feature/template"
	"service/internal/server"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(app *conf.App, serverConf *conf.Server, dataConf *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet, // run grpc/http
		data.ProviderSet,   // DB, Redis, etc
		// modules
		template.ProviderSet, // Template module
		broker.ProviderSet,   // MQTT broker
		// final build kratos.App
		newApp, // final build kratos.App
	))
}
