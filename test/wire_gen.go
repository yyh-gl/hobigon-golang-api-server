// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package test

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/dao"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/rest"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
	"github.com/yyh-gl/hobigon-golang-api-server/cmd/rest/di"
)

// Injectors from wire.go:

func InitTestApp() *di.ContainerAPI {
	gormDB := db.NewDB()
	blog := dao.NewBlog(gormDB)
	slack := dao.NewSlack()
	usecaseBlog := usecase.NewBlog(blog, slack)
	restBlog := rest.NewBlog(usecaseBlog)
	calendar := rest.NewCalendar()
	task := dao.NewTask()
	notification := usecase.NewNotification(task, slack)
	restNotification := rest.NewNotification(notification)
	containerAPI := &di.ContainerAPI{
		HandlerBlog:         restBlog,
		HandlerCalendar:     calendar,
		HandlerNotification: restNotification,
		DB:                  gormDB,
	}
	return containerAPI
}

// wire.go:

var testAppSet = wire.NewSet(infra.APISet, usecase.APISet, rest.WireSet)
