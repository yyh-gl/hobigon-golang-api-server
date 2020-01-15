package main

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra"
	"github.com/yyh-gl/hobigon-golang-api-server/app/interface/rest"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
	"github.com/yyh-gl/hobigon-golang-api-server/cmd/api/di"
)

// TODO: infra, usecaseもapiとcliで分ける
var appSet = wire.NewSet(
	app.APISet,
	infra.WireSet,
	usecase.WireSet,
	rest.WireSet,
)

func initApp() *di.ContainerAPI {
	wire.Build(
		wire.Struct(new(di.ContainerAPI), "*"),
		appSet,
	)
	return nil
}
