//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra"
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
	"github.com/yyh-gl/hobigon-golang-api-server/cmd/rest/di"
)

var appSet = wire.NewSet(
	infra.CLISet,
	usecase.CLISet,
	cli.WireSet,
)

func initApp() *di.ContainerCLI {
	wire.Build(
		wire.Struct(new(di.ContainerCLI), "*"),
		appSet,
	)
	return nil
}
