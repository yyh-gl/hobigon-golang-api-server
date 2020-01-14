package di

import (
	"log"

	"github.com/yyh-gl/hobigon-golang-api-server/handler/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/handler/rest"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/db"
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
