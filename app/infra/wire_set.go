package infra

import (
	"github.com/google/wire"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/dao"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/line"
)

// APISet : infra層のWireSet（API用）
var APISet = wire.NewSet(
	db.NewDB,
	dao.NewBlog,
	dao.NewSlack,
	dao.NewTask,
	line.NewLINE,
	line.NewLINEMessageConfig,
)

// CLISet : infra層のWireSet（CLI用）
var CLISet = wire.NewSet(
	dao.NewSlack,
	dao.NewTask,
	line.NewLINE,
	line.NewLINEMessageConfig,
)
