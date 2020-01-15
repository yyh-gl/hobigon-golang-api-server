package main

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
	"github.com/yyh-gl/hobigon-golang-api-server/cmd/api/di"
)

// TODO: infra, usecaseもapiとcliで分ける
var appSet = wire.NewSet(
	app.CLISet,
	infra.WireSet,
	usecase.WireSet,
	cli.WireSet,
)

func initApp() *di.ContainerCLI {
	wire.Build(
		wire.Struct(new(di.ContainerCLI), "*"),
		appSet,
	)
	return nil
}
