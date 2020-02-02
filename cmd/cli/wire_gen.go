// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/dao"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/service"
	"github.com/yyh-gl/hobigon-golang-api-server/app/interface/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
	"github.com/yyh-gl/hobigon-golang-api-server/cmd/api/di"
)

// Injectors from wire.go:

func initApp() *di.ContainerCLI {
	task := dao.NewTask()
	slack := dao.NewSlack()
	gormDB := db.NewDB()
	birthday := dao.NewBirthday(gormDB)
	ranking := service.NewRanking()
	notification := usecase.NewNotification(task, slack, birthday, ranking)
	cliNotification := cli.NewNotification(notification)
	logger := app.NewCLILogger()
	containerCLI := &di.ContainerCLI{
		HandlerNotification: cliNotification,
		DB:                  gormDB,
		Logger:              logger,
	}
	return containerCLI
}

// wire.go:

var appSet = wire.NewSet(app.CLISet, infra.CLISet, usecase.CLISet, cli.WireSet)
