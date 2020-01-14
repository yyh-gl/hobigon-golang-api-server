package infra

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/igateway"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/irepository"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/iservice"
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
