package main

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/handler/rest"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/igateway"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/irepository"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/iservice"
	"github.com/yyh-gl/hobigon-golang-api-server/usecase"
)

func initBlogHandler() rest.BlogHandler {
	wire.Build(
		irepository.NewBlogRepository,
		igateway.NewSlackGateway,
		usecase.NewBlogUseCase,
		rest.NewBlogHandler,
	)
	return nil
}

func initBirthdayHandler() rest.BirthdayHandler {
	wire.Build(
		irepository.NewBirthdayRepository,
		usecase.NewBirthdayUseCase,
		rest.NewBirthdayHandler,
	)
	return nil
}

func initNotificationHandler() rest.NotificationHandler {
	wire.Build(
		igateway.NewTaskGateway,
		igateway.NewSlackGateway,
		irepository.NewBirthdayRepository,
		iservice.NewNotificationService,
		iservice.NewRankingService,
		usecase.NewNotificationUseCase,
		rest.NewNotificationHandler,
	)
	return nil
}
