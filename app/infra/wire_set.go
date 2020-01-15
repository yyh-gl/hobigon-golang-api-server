package infra

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/igateway"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/irepository"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/iservice"
)

var WireSet = wire.NewSet(
	db.NewDB,
	irepository.NewBlogRepository,
	irepository.NewBirthdayRepository,
	igateway.NewSlackGateway,
	igateway.NewTaskGateway,
	iservice.NewNotificationService,
	iservice.NewRankingService,
)
