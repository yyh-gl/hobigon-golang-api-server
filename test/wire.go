// +build wireinject

package test

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra"
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/http"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
	"github.com/yyh-gl/hobigon-golang-api-server/cmd/http/di"
)

var testAppSet = wire.NewSet(
	app.APISet,
	infra.APISet,
	usecase.APISet,
	http.WireSet,
)

func InitTestApp() *di.ContainerAPI {
	wire.Build(
		wire.Struct(new(di.ContainerAPI), "*"),
		testAppSet,
	)
	return nil
}
