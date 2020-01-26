package infra

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/igateway"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/irepository"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/iservice"
)

// APISet : infra層のWireSet（API用）
var APISet = wire.NewSet(
	db.NewDB,
	irepository.NewBlog,
	irepository.NewBirthday,
	igateway.NewSlack,
	igateway.NewTask,
	iservice.NewNotification,
	iservice.NewRanking,
)

// CLISet : infra層のWireSet（CLI用）
var CLISet = wire.NewSet(
	db.NewDB,
	irepository.NewBirthday,
	igateway.NewSlack,
	igateway.NewTask,
	iservice.NewNotification,
	iservice.NewRanking,
)
