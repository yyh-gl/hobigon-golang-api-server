package di

import (
	"log"

	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/app/interface/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app/interface/rest"
)

type ContainerAPI struct {
	HandlerBlog         rest.BlogHandler
	HandlerBirthday     rest.BirthdayHandler
	HandlerNotification rest.NotificationHandler

	DB     *db.DB
	Logger *log.Logger
}

type ContainerCLI struct {
	HandlerNotification cli.NotificationHandler

	DB     *db.DB
	Logger *log.Logger
}
