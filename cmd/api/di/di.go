package di

import (
	"log"

	"github.com/yyh-gl/hobigon-golang-api-server/handler/rest"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/db"
)

type Container struct {
	HandlerBlog         rest.BlogHandler
	HandlerBirthday     rest.BirthdayHandler
	HandlerNotification rest.NotificationHandler

	DB     *db.DB
	Logger *log.Logger
}
