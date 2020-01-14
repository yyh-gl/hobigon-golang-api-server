package di

import (
	"github.com/yyh-gl/hobigon-golang-api-server/handler/rest"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/db"
)

type Container struct {
	HandlerBlog         rest.BlogHandler
	HandlerBirthday     rest.BirthdayHandler
	HandlerNotification rest.NotificationHandler

	DB *db.DB
}
