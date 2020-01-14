package main

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/cmd/api/di"
	"github.com/yyh-gl/hobigon-golang-api-server/handler/rest"
	"github.com/yyh-gl/hobigon-golang-api-server/infra"
	"github.com/yyh-gl/hobigon-golang-api-server/usecase"
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
