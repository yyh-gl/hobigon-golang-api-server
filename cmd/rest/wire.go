//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra"
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/rest"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
	"github.com/yyh-gl/hobigon-golang-api-server/cmd/rest/di"
)

var appSet = wire.NewSet(
	infra.APISet,
	usecase.APISet,
	rest.WireSet,
)

func initApp() *di.ContainerAPI {
	wire.Build(
		wire.Struct(new(di.ContainerAPI), "*"),
		appSet,
	)
	return nil
}
