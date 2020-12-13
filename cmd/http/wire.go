// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra"
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/http"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
	"github.com/yyh-gl/hobigon-golang-api-server/cmd/http/di"
)

var appSet = wire.NewSet(
	app.APISet,
	infra.APISet,
	usecase.APISet,
	http.WireSet,
)

func initApp() *di.ContainerAPI {
	wire.Build(
		wire.Struct(new(di.ContainerAPI), "*"),
		appSet,
	)
	return nil
}
