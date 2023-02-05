//go:build wireinject
// +build wireinject

package test

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra"
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/rest"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
	"github.com/yyh-gl/hobigon-golang-api-server/cmd/rest/di"
)

var testAppSet = wire.NewSet(
	infra.APISet,
	usecase.APISet,
	rest.WireSet,
)

func InitTestApp() *di.ContainerAPI {
	wire.Build(
		wire.Struct(new(di.ContainerAPI), "*"),
		testAppSet,
	)
	return nil
}
