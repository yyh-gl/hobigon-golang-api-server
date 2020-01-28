// +build wireinject

package test

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra"
	"github.com/yyh-gl/hobigon-golang-api-server/app/interface/rest"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
	"github.com/yyh-gl/hobigon-golang-api-server/cmd/api/di"
)

var testAppSet = wire.NewSet(
	app.APISet,
	infra.APISet,
	usecase.APISet,
	rest.WireSet,
)

func initTestApp() *di.ContainerAPI {
	wire.Build(
		wire.Struct(new(di.ContainerAPI), "*"),
		testAppSet,
	)
	return nil
}
