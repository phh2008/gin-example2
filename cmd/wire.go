//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"com.example/example/bootstrap"
	"com.example/example/pkg"
	"com.example/example/pkg/config"
	"com.example/example/repository"
	"com.example/example/service"
	"com.example/example/web/controller"
	"com.example/example/web/middleware"
	"com.example/example/web/router"
	"github.com/google/wire"
)

func BuildServer(conf *config.Config) *bootstrap.Server {
	wire.Build(
		pkg.ProviderSet,
		repository.ProviderSet,
		service.ProviderSet,
		controller.ProviderSet,
		middleware.ProviderSet,
		router.ProviderSet,
		bootstrap.ServerSet,
	)
	return new(bootstrap.Server)
}
