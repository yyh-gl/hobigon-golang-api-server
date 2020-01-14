package main

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/cmd/api/di"
	"github.com/yyh-gl/hobigon-golang-api-server/handler/rest"
	"github.com/yyh-gl/hobigon-golang-api-server/infra"
	"github.com/yyh-gl/hobigon-golang-api-server/usecase"
)

var appSet = wire.NewSet(
	infra.WireSet,
	usecase.WireSet,
	rest.WireSet,
)

func initApp() *di.Container {
	wire.Build(
		wire.Struct(new(di.Container), "*"),
		appSet,
	)
	return nil
}
